package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/karldoenitz/Tigo/TigoWeb"
	"github.com/karldoenitz/tission/session/utils"
	"github.com/patrickmn/go-cache"
	"reflect"
	"time"
)

var (
	si SessionInterface
)

type SessionInterface struct {
	Expire time.Duration
}

func (sip *SessionInterface) initCache() {
	cacheManager = produceCacheManager(sip.Expire, sip.Expire)
}

func (sip *SessionInterface) NewSessionManager() TigoWeb.SessionManager {
	// 这里初始化缓存
	return &SessionManager{expire: sip.Expire}
}

func (sip *SessionInterface) GetCache() *cache.Cache {
	if cacheManager != nil {
		return cacheManager
	}
	sip.initCache()
	return cacheManager
}

type SessionManager struct {
	expire time.Duration
}

func (sm *SessionManager) GenerateSession(expire int) TigoWeb.Session {
	session := Session{}
	session.sessionId = utils.GetSessionId()
	sm.expire = time.Duration(expire)
	session.expire = sm.expire
	Set(session.sessionId, "", sm.expire)
	return &session
}

func (sm *SessionManager) GetSessionBySid(sid string) TigoWeb.Session {
	session := Session{sessionId: sid}
	// 从缓存中获取
	_, isFound := Get(sid)
	if !isFound {
		return &session
	}
	session.sessionId = sid
	session.expire = sm.expire
	return &session
}

func (sm *SessionManager) DeleteSession(sid string) {
	Del(sid)
}

type Session struct {
	sessionId string
	expire    time.Duration
}

func (s *Session) updateSession() (err error) {
	// 向缓存中设置
	if err := Set(s.sessionId, "", s.expire); err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (s *Session) Get(key string, value interface{}) (err error) {
	defer func() {
		fatalError := recover()
		if fatalError != nil {
			err = errors.New(fmt.Sprintf("session值类型与返回值不匹配: error(%#v)", fatalError))
			return
		}
	}()
	sessionKey := fmt.Sprintf("tission_%s_%s", s.sessionId, key)
	sv, isExisted := Get(sessionKey)
	if !isExisted {
		return errors.New(fmt.Sprintf("session value of key(%s) is nil", key))
	}
	valPtr := reflect.ValueOf(value).Elem()
	switch valPtr.Type().Kind() {
	case reflect.String:
		v := sv.(string)
		valPtr.SetString(v)
	case reflect.Bool:
		v := sv.(bool)
		valPtr.SetBool(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v := sv.(uint64)
		valPtr.SetUint(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := sv.(int64)
		valPtr.SetInt(v)
	case reflect.Float32, reflect.Float64:
		v := sv.(float64)
		valPtr.SetFloat(v)
	case reflect.Map, reflect.Slice, reflect.Struct:
		b, e := json.Marshal(sv)
		if e != nil {
			return e
		}
		if e := json.Unmarshal(b, value); e != nil {
			return e
		}
	}
	return
}

func (s *Session) Set(key string, value interface{}) (err error) {
	sessionKey := fmt.Sprintf("tission_%s_%s", s.sessionId, key)
	if err := Set(sessionKey, value, s.expire); err != nil {
		fmt.Printf("set %s %v to session failed => %s", key, value, err.Error())
	}
	return s.updateSession()
}

func (s *Session) Delete(key string) {
	sessionKey := fmt.Sprintf("tission_%s_%s", s.sessionId, key)
	Del(sessionKey)
	_ = s.updateSession()
}

func (s *Session) SessionId() (sid string) {
	return s.sessionId
}

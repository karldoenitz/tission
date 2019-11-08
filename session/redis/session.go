package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/karldoenitz/Tigo/TigoWeb"
	"reflect"
	"time"
)

var si SessionInterface

var (
	IP      string
	Port    string
	MaxIdle int
	Timeout int
	Pwd     string
	Auth    string
	DbNo    int
)

type SessionInterface struct {
	IP      string
	Port    string
	MaxIdle int
	Timeout int
	Pwd     string
	Auth    string
	DbNo    int
	Expire  int64
}

func (sip *SessionInterface) initRedisPool() {
	addr := fmt.Sprintf("%s:%s", IP, Port)
	redisPool = produceRedisPool(addr, MaxIdle, Timeout, Auth, DbNo, Pwd)
}

func (sip *SessionInterface) GetRedisPool() *redis.Pool {
	if redisPool != nil {
		return redisPool
	}
	sip.initRedisPool()
	return redisPool
}

func (sip *SessionInterface) NewSessionManager() TigoWeb.SessionManager {
	IP = sip.IP
	Port = sip.Port
	MaxIdle = sip.MaxIdle
	Timeout = sip.Timeout
	Pwd = sip.Pwd
	Auth = sip.Auth
	DbNo = sip.DbNo
	sip.initRedisPool()
	return &SessionManager{expire: sip.Expire}
}

type SessionManager struct {
	expire int64
}

func (sm *SessionManager) GenerateSession(expire int) TigoWeb.Session {
	session := Session{}
	session.sessionId = getSessionId()
	session.value = make(map[string]interface{})
	if expire > 0 {
		sm.expire = int64(expire) * int64(time.Second)
	}
	session.expire = sm.expire
	Set(session.sessionId, session.value, time.Duration(sm.expire))
	return &session
}

func (sm *SessionManager) GetSessionBySid(sid string) TigoWeb.Session {
	session := Session{}
	value, isFound := Get(sid)
	if !isFound {
		return &session
	}
	session.sessionId = sid
	session.value = make(map[string]interface{})
	session.expire = sm.expire
	if session.expire == 0 {
		session.expire = int64(3600) * int64(time.Second)
	}
	json.Unmarshal(value, &(session.value))
	return &session
}

func (sm *SessionManager) DeleteSession(sid string) {
	Del(sid)
}

type Session struct {
	value     map[string]interface{}
	sessionId string
	expire    int64
}

func (s *Session) updateSession() (err error) {
	data, err := json.Marshal(s.value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err, _ := Set(s.sessionId, data, time.Duration(s.expire)); err != nil {
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
	sv, isExisted := s.value[key]
	if !isExisted {
		return errors.New(fmt.Sprintf("session value of key(%s) is nil", key))
	}
	valPtr := reflect.ValueOf(value).Elem()
	switch valPtr.Kind() {
	case reflect.String:
		v := sv.(string)
		valPtr.SetString(v)
	case reflect.Bool:
		v := sv.(bool)
		valPtr.SetBool(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v := sv.(float64)
		valPtr.SetUint(uint64(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := sv.(float64)
		valPtr.SetInt(int64(v))
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
	s.value[key] = value
	return s.updateSession()
}

func (s *Session) Delete(key string) {
	if _, isExisted := s.value[key]; isExisted {
		delete(s.value, key)
	}
	s.updateSession()
}

func (s *Session) SessionId() (sid string) {
	return s.sessionId
}

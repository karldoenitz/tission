package redis

import (
	"encoding/json"
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
	return &SessionManager{}
}

type SessionManager struct {
	expire int64
}

func (sm *SessionManager) GenerateSession(expire int) TigoWeb.Session {
	sessionId := fmt.Sprintf("%d", time.Now().Local().Unix())
	session := Session{}
	session.sessionId = sessionId
	session.value = make(map[string]interface{})
	sm.expire = int64(expire) * int64(time.Second)
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

}

type Session struct {
	value     map[string]interface{}
	sessionId string
	expire    int64
}

func (s *Session) updateSession() {
	data, _ := json.Marshal(s.value)
	if err, _ := Set(s.sessionId, data, time.Duration(s.expire)); err != nil {
		fmt.Println(err.Error())
	}
}

func (s *Session) Get(key string, value interface{}) (err error) {
	v := s.value[key].(string)
	reflect.ValueOf(value).Elem().SetString(v)
	return
}

func (s *Session) Set(key string, value interface{}) (err error) {
	s.value[key] = value
	s.updateSession()
	return
}

func (s *Session) Delete(key string) {

}

func (s *Session) SessionId() (sid string) {
	return s.sessionId
}

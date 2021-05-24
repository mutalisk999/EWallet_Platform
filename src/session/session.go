package session

import (
	"github.com/kataras/iris/core/errors"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

type SessionValue struct {
	AcctId     int
	Role       int
	CellNumber string
	RealName   string
	IdCard     string
	PubKey     string
	CreateTime time.Time
	UpdateTime time.Time
}

type SessionMgr struct {
	SessionStore map[string]SessionValue
	Mutex        *sync.Mutex
}

var GlobalSessionMgr *SessionMgr

func InitSessionMgr() {
	GlobalSessionMgr = new(SessionMgr)
	GlobalSessionMgr.InitSessionStore()
	GlobalSessionMgr.Mutex = new(sync.Mutex)
}

func (m *SessionMgr) InitSessionStore() {
	m.SessionStore = make(map[string]SessionValue)
}

func (m SessionMgr) HasSessionId(sessionId string) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	_, ok := m.SessionStore[sessionId]
	return ok
}

func (m SessionMgr) GetSessionValue(sessionId string) (SessionValue, bool) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	sessionValue, ok := m.SessionStore[sessionId]
	if !ok {
		return SessionValue{}, false
	}
	sessionValue.UpdateTime = time.Now()
	m.SessionStore[sessionId] = sessionValue
	return sessionValue, true
}

func (m SessionMgr) DeleteSessionValueByAcctId(acctId int) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var sessionIds []string
	for k, v := range m.SessionStore {
		if v.AcctId == acctId {
			sessionIds = append(sessionIds, k)
		}
	}
	if len(sessionIds) == 0 {
		return false
	}

	for _, sessionId := range sessionIds {
		delete(m.SessionStore, sessionId)
	}
	return true
}

//func (m SessionMgr) RefreshSessionValue(sessionId string) {
//	m.Mutex.Lock()
//	defer m.Mutex.Unlock()
//
//	sessionValue, ok := m.SessionStore[sessionId]
//	if ok {
//		sessionValue.UpdateTime = time.Now()
//		m.SessionStore[sessionId] = sessionValue
//	}
//}

func (m *SessionMgr) NewSessionValue(sessionValue SessionValue) (string, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	tryCount := 0
	for {
		if tryCount >= 10 {
			return "", errors.New("can not create valid session id")
		}
		u, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		_, ok := m.SessionStore[u.String()]
		if !ok {
			sessionValue.CreateTime = time.Now()
			sessionValue.UpdateTime = time.Now()
			m.SessionStore[u.String()] = sessionValue
			return u.String(), nil
		}
		tryCount += 1
	}
}

func (m *SessionMgr) DeleteSessionValue(sessionId string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.SessionStore, sessionId)
}

func (m SessionMgr) IsAdmin(sessionId string) (bool, error) {
	sessionValue, ok := m.GetSessionValue(sessionId)
	if !ok {
		return false, errors.New("session id not exist")
	}
	if sessionValue.Role == 0 {
		return true, nil
	}
	return false, nil
}

func (m SessionMgr) IsAccountant(sessionId string) (bool, error) {
	sessionValue, ok := m.GetSessionValue(sessionId)
	if !ok {
		return false, errors.New("session id not exist")
	}
	if sessionValue.Role == 1 {
		return true, nil
	}
	return false, nil
}

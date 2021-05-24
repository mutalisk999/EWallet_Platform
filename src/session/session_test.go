package session

import (
	"fmt"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	InitSessionMgr()
	sessionValue := SessionValue{0, 0, "1234567", "张三", "12345678", "1234567890", time.Now(), time.Now()}
	sid, _ := GlobalSessionMgr.NewSessionValue(sessionValue)
	fmt.Println("sid", sid)

	sessionValue, ok := GlobalSessionMgr.GetSessionValue(sid)
	if ok {
		fmt.Println("sessionValue:", sessionValue)
	}

	isAdmin, _ := GlobalSessionMgr.IsAdmin(sid)
	fmt.Println("isAdmin:", isAdmin)

	isAccountant, _ := GlobalSessionMgr.IsAccountant(sid)
	fmt.Println("isAccountant:", isAccountant)

	GlobalSessionMgr.DeleteSessionValue(sid)

	hasSessionId := GlobalSessionMgr.HasSessionId(sid)
	fmt.Println("hasSessionId:", hasSessionId)
}

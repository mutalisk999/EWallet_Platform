package model

import (
	"sync"
	"time"
)

type tblOperatorLog struct {
	Logid      int       `xorm:"pk INTEGER autoincr"`
	Acctid     int       `xorm:"INT NOT NULL"`
	Optype     int       `xorm:"INT NOT NULL"`
	Content    string    `xorm:"TEXT NOT NULL"`
	Createtime time.Time `xorm:"created"`
}

type tblOperationLogMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblOperationLogMgr) Init() {
	t.TableName = "tbl_operator_log"
	t.Mutex = new(sync.Mutex)
}

func (t *tblOperationLogMgr) GetOperatorLogs(acctId []int, opType []int, opTime [2]string, offSet int, limit int) (int, []tblOperatorLog, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	dbSession := GetDBEngine().Where("")
	if acctId != nil && len(acctId) != 0 {
		dbSession = dbSession.In("acctid", acctId)
	}
	if opType != nil && len(opType) != 0 {
		dbSession = dbSession.In("optype", opType)
	}
	if opTime[0] != "" {
		dbSession = dbSession.And("createtime > ?", opTime[0])
	}
	if opTime[1] != "" {
		dbSession = dbSession.And("createtime < ?", opTime[1])
	}
	var log tblOperatorLog
	total, err := dbSession.Count(&log)
	if err != nil {
		return 0, nil, err
	}

	dbSession2 := GetDBEngine().Where("")
	if acctId != nil && len(acctId) != 0 {
		dbSession2 = dbSession2.In("acctid", acctId)
	}
	if opType != nil && len(opType) != 0 {
		dbSession2 = dbSession2.In("optype", opType)
	}
	if opTime[0] != "" {
		dbSession2 = dbSession2.And("createtime > ?", opTime[0])
	}
	if opTime[1] != "" {
		dbSession2 = dbSession2.And("createtime < ?", opTime[1])
	}
	opLogs := make([]tblOperatorLog, 0)
	dbSession2.Limit(limit, offSet).Desc("createtime").Find(&opLogs)
	return int(total), opLogs, nil
}

func (t *tblOperationLogMgr) NewOperatorLog(acctId int, opType int, content string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var log tblOperatorLog
	log.Acctid = acctId
	log.Optype = opType
	log.Content = content
	log.Createtime = time.Now()
	_, err := GetDBEngine().Insert(&log)
	if err != nil {
		return 0, err
	}
	return log.Logid, nil
}

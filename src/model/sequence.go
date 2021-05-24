package model

import (
	"fmt"
	"sync"
	"time"
)

type tblSequence struct {
	Seqvalue   int       `xorm:"pk INTEGER autoincr"`
	Idtype     int       `xorm:"INT NOT NULL"`
	State      int       `xorm:"INT NOT NULL"`
	Createtime time.Time `xorm:"created"`
	Updatetime time.Time `xorm:"DATETIME"`
}

type tblSequenceMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblSequenceMgr) Init() {
	t.TableName = "tbl_sequence"
	t.Mutex = new(sync.Mutex)
}

func (t *tblSequenceMgr) NewSequence(idType int) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var sequence tblSequence
	sequence.Idtype = idType
	sequence.State = 0
	sequence.Createtime = time.Now()
	sequence.Updatetime = time.Now()
	_, err := GetDBEngine().Insert(&sequence)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return sequence.Seqvalue, nil
}

func (t *tblSequenceMgr) QuerySequence(idType int, seqValue int) (bool, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var sequence tblSequence
	sequence.Seqvalue = seqValue
	sequence.Idtype = idType
	sequence.State = 0
	total, err := GetDBEngine().Count(&sequence)
	if err != nil {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	return true, nil
}

func (t *tblSequenceMgr) VerifySequence(idType int, seqValue int) (bool, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var sequence tblSequence
	sequence.Seqvalue = seqValue
	sequence.Idtype = idType
	sequence.State = 0
	total, err := GetDBEngine().Delete(&sequence)
	if err != nil {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	return true, nil
}

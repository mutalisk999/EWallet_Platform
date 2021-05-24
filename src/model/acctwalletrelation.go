package model

import (
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
	"fmt"
)

type tblAcctWalletRelation struct {
	Relationid int       `xorm:"pk INTEGER autoincr NOT NULL"`
	Acctid     int       `xorm:"INTEGER UNIQUE(walletid) NOT NULL"`
	Walletid   int       `xorm:"INTEGER NOT NULL"`
	Createtime time.Time `xorm:"DATETIME"`
}
type tblAcctWalletRelationMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblAcctWalletRelationMgr) Init() {
	t.TableName = "tbl_acct_wallet_relation"
	t.Mutex = new(sync.Mutex)
}
func (t *tblAcctWalletRelationMgr) InsertRelation(accid int, walletid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var relation tblAcctWalletRelation
	result, err := GetDBEngine().Where("acctid=? and walletid=?", accid,walletid).Get(&relation)
	if result{
		return nil
	}
	if err!=nil{
		fmt.Println(err.Error())
		return err
	}
	relation.Createtime = time.Now()
	relation.Acctid = accid
	relation.Walletid = walletid
	_, err = GetDBEngine().Insert(&relation)
	return err
}
func (t *tblAcctWalletRelationMgr) GetRelationByRelationId(rid int) (tblAcctWalletRelation, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var relation tblAcctWalletRelation
	result, err := GetDBEngine().Where("relationid=?", rid).Get(&relation)
	if err != nil {
		return relation, err
	}
	if result {
		return relation, nil
	}
	return relation, errors.New("no find relation")
}


func (t *tblAcctWalletRelationMgr) GetRelationsByWalletId(walletid int) ([]tblAcctWalletRelation, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	relations := make([]tblAcctWalletRelation, 0)
	err := GetDBEngine().Where("walletid=?", walletid).Find(&relations)
	if err != nil {
		return relations, err
	}
	return relations, nil
}
func (t *tblAcctWalletRelationMgr) GetRelationsByAcctId(Accid int) ([]tblAcctWalletRelation, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	relations := make([]tblAcctWalletRelation, 0)
	err := GetDBEngine().Where("acctid=?", Accid).Find(&relations)
	if err != nil {
		return relations, err
	}
	return relations, nil
}
func (t *tblAcctWalletRelationMgr) DeleteRelation(rid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var relation tblAcctWalletRelation
	relation.Relationid = rid
	_, err := GetDBEngine().Delete(relation)
	return err
}
func (t *tblAcctWalletRelationMgr) DeleteRelationByWalletId(wid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var relation tblAcctWalletRelation
	relation.Walletid = wid
	_, err := GetDBEngine().Delete(relation)
	return err
}
func (t *tblAcctWalletRelationMgr) DeleteRelationByAccountId(aid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var relation tblAcctWalletRelation
	relation.Acctid = aid
	_, err := GetDBEngine().Delete(relation)
	return err
}

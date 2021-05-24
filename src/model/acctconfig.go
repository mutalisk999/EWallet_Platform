package model

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
)

type tblAcctConfig struct {
	Acctid     int       `xorm:"pk INTEGER autoincr"`
	Cellphone  string    `xorm:"VARCHAR(64) NOT NULL"`
	Realname   string    `xorm:"VARCHAR(64) NOT NULL"`
	Idcard     string    `xorm:"VARCHAR(64) NOT NULL"`
	Pubkey     string    `xorm:"VARCHAR(512) NOT NULL UNIQUE"`
	Role       int       `xorm:"INT NOT NULL"`
	State      int       `xorm:"INT NOT NULL"`
	Createtime time.Time `xorm:"DATETIME"`
	Updatetime time.Time `xorm:"DATETIME"`
}

type tblAcctConfigMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblAcctConfigMgr) Init() {
	t.TableName = "tbl_acct_config"
	t.Mutex = new(sync.Mutex)
}

func (t *tblAcctConfigMgr) VerifyUnique(CellPhone string, RealName string, IdCard string, Pubkey string) error {
	t.Mutex.Lock()
	var acct tblAcctConfig
	state := []int{0,1,2}
	result, err := GetDBEngine().In("state",state).Where("cellphone=? or realname=? or idcard=? or pubkey=?", CellPhone,RealName,IdCard,Pubkey).Get(&acct)
	if result {
		t.Mutex.Unlock()
		return errors.New("key already exist!")
	}
	t.Mutex.Unlock()
	return err
}

func (t *tblAcctConfigMgr) InsertAcct(CellPhone string, RealName string, IdCard string, Pubkey string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	count, err := GetDBEngine().Table(t.TableName).Count(&acct)
	acct.Cellphone = CellPhone
	acct.Realname = RealName
	acct.Idcard = IdCard
	acct.Pubkey = Pubkey
	acct.State = 0
	acct.Createtime = time.Now()
	acct.Updatetime = time.Now()
	if err != nil {
		fmt.Println("account count ",err.Error())
	}
	if count > 0 {
		acct.Role = 1
	} else {
		acct.Role = 0
		acct.State = 1
	}

	_, err = GetDBEngine().Insert(&acct)
	return err
}

func (t *tblAcctConfigMgr) GetAdminId() int {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := GetDBEngine().Where("role=0").Get(&acct)
	if err != nil {
		return -1
	}
	if !result {
		return -1
	}
	return acct.Acctid
}

func (t *tblAcctConfigMgr) UpdateAcct(acctid int, CellPhone string, RealName string, IdCard string, Pubkey string, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := GetDBEngine().Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.Cellphone = CellPhone
		acct.Idcard = IdCard
		acct.Pubkey = Pubkey
		acct.State = State
		acct.Updatetime = time.Now()
		_, err = GetDBEngine().Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) UpdateAcctState(acctid int, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := GetDBEngine().Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.State = State
		acct.Updatetime = time.Now()
		_, err = GetDBEngine().Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) ActiveAcct(acctid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := GetDBEngine().Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.State = 1
		acct.Updatetime = time.Now()
		_, err = GetDBEngine().Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) FreezeAcct(acctid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := GetDBEngine().Where("acctid=?", acctid).Get(&acct)
	if err != nil {
		return err
	}
	if result {
		acct.State = 2
		acct.Updatetime = time.Now()
		_, err = GetDBEngine().Where("acctid=?", acctid).Update(&acct)
		return err
	} else {
		return errors.New("not find account")
	}
}

func (t *tblAcctConfigMgr) DeleteAcct(acctid int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	acct.Acctid = acctid
	_, err := GetDBEngine().Delete(acct)
	return err
}

func (t *tblAcctConfigMgr) FindAccountByPubkey(pubkey string) (tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	result, err := GetDBEngine().Where("pubkey=?", pubkey).Get(&acct)
	if result {
		return acct, err
	} else {
		return acct, errors.New("Not Found")
	}

}

func (t *tblAcctConfigMgr) ListNormalAccount(state []int, limit, offset int) ([]tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct []tblAcctConfig
	err := GetDBEngine().In("state", state).And("role=?", 1).Limit(limit, offset).Find(&acct)

	return acct, err
}

func (t *tblAcctConfigMgr) GetNormalAccountCount(state []int) (int64, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	count, err := GetDBEngine().In("state", state).And("role=?", 1).Count(&acct)

	return count, err
}

func (t *tblAcctConfigMgr) GetAccountById(acctId int) (tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	exist, err := GetDBEngine().Where("acctid=?", acctId).Get(&acct)
	if err != nil {
		return acct, err
	}
	if !exist {

		return acct, errors.New("Account Not Found")
	}

	return acct, err
}



func (t *tblAcctConfigMgr) GetAccountCount() (int64, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	count, err := GetDBEngine().Count(&acct)
	return count,err
}

func (t *tblAcctConfigMgr) GetAccountIdByPubkey(pubkey string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var acct tblAcctConfig
	exist, err := GetDBEngine().Where("pubkey=?", pubkey).Get(&acct)
	if err != nil {
		return -1, err
	}
	if !exist {

		return -1, errors.New("Account Not Found")
	}

	return acct.Acctid, err
}

func (t *tblAcctConfigMgr) GetAccountsByIds(acctIds []int) ([]tblAcctConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var accts []tblAcctConfig
	if len(acctIds)>0{
		err := GetDBEngine().In("acctid", acctIds).Find(&accts)
		if err != nil {
			return nil, err
		}
	}else if len(acctIds)==0{
		err := GetDBEngine().Find(&accts)
		if err != nil {
			return nil, err
		}
	}

	return accts, nil
}

package model

import (
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
)

type tblPubkeyPool struct {
	Keyindex   int       `xorm:"INT pk NOT NULL"`
	Pubkey     string    `xorm:"VARCHAR(256) NOT NULL UNIQUE"`
	Isused     bool      `xorm:"BOOL NOT NULL"`
	Createtime time.Time `xorm:"DATETIME"`
	Usedtime   time.Time `xorm:"DATETIME"`
}

type tblPubKeyPoolMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblPubKeyPoolMgr) Init() {
	t.TableName = "tbl_pubkey_pool"
	t.Mutex = new(sync.Mutex)
}

func (t *tblPubKeyPoolMgr) InsertPubkey(KeyIndex int, Pubkey string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	count, err := GetDBEngine().Where("keyindex=?", KeyIndex).Count(&pubkey_record)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("key already exist!")
	}
	pubkey_record.Keyindex = KeyIndex
	pubkey_record.Pubkey = Pubkey
	pubkey_record.Isused = false
	pubkey_record.Createtime = time.Now()
	pubkey_record.Usedtime = time.Now()

	_, err = GetDBEngine().Insert(&pubkey_record)
	return err
}

func (t *tblPubKeyPoolMgr) UpdatePubkey(KeyIndex int, Pubkey string, isuse bool) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	exist, err := GetDBEngine().Where("keyindex=?", KeyIndex).Get(&pubkey_record)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("key not found!")
	}
	pubkey_record.Keyindex = KeyIndex
	pubkey_record.Pubkey = Pubkey
	pubkey_record.Isused = isuse
	pubkey_record.Usedtime = time.Now()

	_, err = GetDBEngine().Where("keyindex=?", KeyIndex).Cols("isused").Update(&pubkey_record)
	return err
}

func (t *tblPubKeyPoolMgr) UsePubkey(KeyIndex int) (string, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var pubkey_record tblPubkeyPool
	exist, err := GetDBEngine().Where("keyindex=?", KeyIndex).Get(&pubkey_record)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("key not found!")
	}
	pubkey_record.Keyindex = KeyIndex
	pubkey_record.Isused = true
	pubkey_record.Usedtime = time.Now()

	_, err = GetDBEngine().Where("keyindex=?", KeyIndex).Cols("isused").Update(&pubkey_record)
	return pubkey_record.Pubkey, err
}
func (t *tblPubKeyPoolMgr) GetAnUnusedKeyIndex() (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var pubkey_record tblPubkeyPool
	exist, err := GetDBEngine().Where("isused=?", false).Get(&pubkey_record)
	if err != nil {
		return -1, errors.New("query key error!")
	}
	if !exist {
		return -1, errors.New("key not found!")
	}
	pubkey_record.Isused = true
	pubkey_record.Usedtime = time.Now()
	_, err = GetDBEngine().Where("keyindex=?", pubkey_record.Keyindex).Update(&pubkey_record)
	return pubkey_record.Keyindex, err

}

func (t *tblPubKeyPoolMgr) QueryPubKeyByKeyIndex(keyIndex int) (string, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var pubkey_record tblPubkeyPool
	exist, err := GetDBEngine().Where("keyindex=?", keyIndex).And("isused=?", true).Get(&pubkey_record)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.New("keyindex not found!")
	}
	return pubkey_record.Pubkey, nil
}

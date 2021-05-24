package model

import (
	"github.com/kataras/iris/core/errors"
	"sync"
	"time"
)

type tblWalletConfig struct {
	Walletid     int       `xorm:"pk INTEGER autoincr"`
	Coinid       int       `xorm:"INTEGER"`
	Walletname   string    `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Keyindex     int       `xorm:"INTEGER NOT NULL UNIQUE index"`
	Address      string    `xorm:"VARCHAR(64) NOT NULL UNIQUE index"`
	Destaddress  string    `xorm:"TEXT"`
	Needsigcount int       `xorm:"INTEGER NOT NULL"`
	Fee          string    `xorm:"VARCHAR(64)"`
	Gasprice     string    `xorm:"VARCHAR(64)"`
	Gaslimit     string    `xorm:"VARCHAR(64)"`
	State        int       `xorm:"INTEGER NOT NULL"`
	Createtime   time.Time `xorm:"DATETIME"`
	Updatetime   time.Time `xorm:"DATETIME"`
}
type tblWalletConfigMgr struct {
	TableName string
	Mutex     *sync.Mutex
}
type WalletRelationinfo struct {
	tblWalletConfig       `xorm:"extends"`
	tblAcctWalletRelation `xorm:"extends"`
}

func (WalletRelationinfo) TableName() string {
	return "tbl_wallet_config"
}
func (t *tblWalletConfigMgr) Init() {
	t.TableName = "tbl_wallet_config"
	t.Mutex = new(sync.Mutex)
}
func (t *tblWalletConfigMgr) InsertWallet(AssetId int, Walletname string, Keyindex int, Address string, Destaddress string, Needsigcount int, Fee string, GasPrice string, GasLimit string, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	wallet.Address = Address
	wallet.Coinid = AssetId
	wallet.Createtime = time.Now()
	wallet.Destaddress = Destaddress
	wallet.Keyindex = Keyindex
	wallet.Needsigcount = Needsigcount
	wallet.Fee = Fee
	wallet.Gasprice = GasPrice
	wallet.Gaslimit = GasLimit
	wallet.State = State
	wallet.Updatetime = time.Now()
	wallet.Walletname = Walletname
	_, err := GetDBEngine().Insert(&wallet)
	return err
}
func (t *tblWalletConfigMgr) UpdateWallet(walletid int, Walletname string, Destaddress string, Needsigcount int, Fee string, GasPrice string, GasLimit string, State int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := GetDBEngine().Where("walletid=?", walletid).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.Destaddress = Destaddress
		wallet.Needsigcount = Needsigcount
		wallet.Fee = Fee
		wallet.Gasprice = GasPrice
		wallet.Gaslimit = GasLimit
		wallet.State = State
		wallet.Updatetime = time.Now()
		wallet.Walletname = Walletname
		_, err := GetDBEngine().Where("walletid=?", walletid).Cols("needsigcount","state","updatetime","walletname","destaddress","fee","gasprice","gaslimit").Update(&wallet)
		return err
	} else {
		return errors.New("not find wallet")
	}

}
func (t *tblWalletConfigMgr) ListWallets(coinids []int, state []int, acctids []int, offset int, limit int) ([]tblWalletConfig, int64, error) {

	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	result := make([]tblWalletConfig, 0)
	walls := make([]tblAcctWalletRelation, 0)
	wallids := make(map[int]int, 0)
	wids:=make([]int, 0)
	if len(acctids)!=0{
		err :=GetDBEngine().In("acctid",acctids).Find(&walls)
		if err != nil {
			return result, 0, err
		}
		for _,rela :=range walls{
			wallids[rela.Walletid]=0
		}
	}
	for k,_:=range wallids{
		wids=append(wids, k)
	}
	if len(acctids)!=0&&len(wids)==0{
		return result,0,nil
	}
	dbSe := GetDBEngine().Where("")
	if len(wids)!=0{
		dbSe = dbSe.In("walletid", wids)
	}
	if len(coinids)!=0{
		dbSe = dbSe.In("coinid", coinids)
	}
	if len(state) != 0 {
		dbSe = dbSe.In("state", state)
	}
	err := dbSe.Limit(limit,offset).Find(&result)
	if err != nil {
		return result, 0, err
	}
	dbSe = GetDBEngine().Where("")
	if len(wids)!=0{
		dbSe = dbSe.In("walletid", wids)
	}
	if len(coinids)!=0{
		dbSe = dbSe.In("coinid", coinids)
	}
	if len(state) != 0 {
		dbSe = dbSe.In("state", state)
	}
	var EmptyWa tblWalletConfig
	total,err:= dbSe.Count(EmptyWa)
	if err != nil {
		return result, 0, err
	}
	return result, total, nil
}
func (t *tblWalletConfigMgr) GetWalletById(id int) (tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := GetDBEngine().Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return wallet, err
	}
	if result {
		return wallet, nil
	}
	return wallet, errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) GetWalletsByIds(ids []int) ([]tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallets []tblWalletConfig
	err := GetDBEngine().In("walletid", ids).Find(&wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (t *tblWalletConfigMgr) GetWalletByName(name string) (tblWalletConfig, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := GetDBEngine().Where("walletname=?", name).Get(&wallet)
	if err != nil {
		return wallet, err
	}
	if result {
		return wallet, nil
	}
	return wallet, errors.New("no find wallet")
}

func (t *tblWalletConfigMgr) ChangeWalletState(id int, sta int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := GetDBEngine().Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.Walletid = id
		wallet.Updatetime = time.Now()
		wallet.State = sta
		_, err := GetDBEngine().Cols("state","updatetime").Update(wallet)
		return err
	}
	return errors.New("no find wallet")
}
func (t *tblWalletConfigMgr) ActiviteWallet(id int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := GetDBEngine().Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.State = 1
		wallet.Updatetime = time.Now()
		_, err := GetDBEngine().Where("walletid=?", id).Update(&wallet)
		return err
	}
	return errors.New("no find wallet")
}
func (t *tblWalletConfigMgr) FreezeWallet(id int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	result, err := GetDBEngine().Where("walletid=?", id).Get(&wallet)
	if err != nil {
		return err
	}
	if result {
		wallet.State = 0
		wallet.Updatetime = time.Now()
		_, err := GetDBEngine().Where("walletid=?", id).Update(&wallet)
		return err
	}
	return errors.New("no find wallet")
}
func (t *tblWalletConfigMgr) DeleteWallet(id int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	var wallet tblWalletConfig
	wallet.Walletid = id
	_, err := GetDBEngine().Delete(wallet)
	return err
}

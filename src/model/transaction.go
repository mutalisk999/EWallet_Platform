package model

import (
	"errors"
	"sync"
	"time"
)

type tblTransaction struct {
	Trxid         int    `xorm:"pk INTEGER autoincr"`
	Rawtrxid      string `xorm:"VARCHAR(128)"`
	Walletid      int    `xorm:"INT NOT NULL"`
	Coinid        int    `xorm:"INT NOT NULL"`
	Contractaddr  string `xorm:"VARCHAR(128)"`
	Acctid        int    `xorm:"INT NOT NULL"`
	Fromaddr      string `xorm:"VARCHAR(128) NOT NULL"`
	Toaddr        string `xorm:"VARCHAR(128) NOT NULL"`
	Amount        string `xorm:"VARCHAR(128) NOT NULL"`
	Feecost       string `xorm:"VARCHAR(128)"`
	Trxtime       time.Time
	Needconfirm   int    `xorm:"INT NOT NULL"`
	Confirmed     int    `xorm:"INT NOT NULL"`
	Acctconfirmed string `xorm:"VARCHAR(1024) NOT NULL"`
	Fee           string `xorm:"VARCHAR(128)"`
	Gasprice      string `xorm:"VARCHAR(128)"`
	Gaslimit      string `xorm:"VARCHAR(128)"`
	State         int    `xorm:"INT NOT NULL"`
	Signature     string `xorm:"VARCHAR(512) NOT NULL"`
}

type tblTransactionMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblTransactionMgr) Init() {
	t.TableName = "tbl_transaction"
	t.Mutex = new(sync.Mutex)
}

func (t *tblTransactionMgr) NewTransaction(walletId int, coinId int, contractAddr string, acctId int, from string, to string,
	amount string, needConfirm int, fee string, gasPrice string, gasLimit string, signature string) (int, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var transaction tblTransaction
	transaction.Walletid = walletId
	transaction.Coinid = coinId
	transaction.Contractaddr = contractAddr
	transaction.Acctid = acctId
	transaction.Fromaddr = from
	transaction.Toaddr = to
	transaction.Amount = amount
	transaction.Trxtime = time.Now()
	transaction.Needconfirm = needConfirm
	transaction.Confirmed = 0
	transaction.Fee = fee
	transaction.Gasprice = gasPrice
	transaction.Gaslimit = gasLimit
	transaction.State = 0
	transaction.Signature = signature
	_, err := GetDBEngine().Insert(&transaction)
	if err != nil {
		return 0, err
	}
	return transaction.Trxid, nil
}

func (t *tblTransactionMgr) GetTransactionById(trxId int) (tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var trx tblTransaction
	result, err := GetDBEngine().Where("trxid=?", trxId).Get(&trx)
	if err != nil {
		return trx, err
	}
	if result {
		return trx, nil
	}
	return trx, errors.New("no find transaction")
}

func (t *tblTransactionMgr) UpdateTransaction(transaction tblTransaction) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	_, err := GetDBEngine().Id(transaction.Trxid).Update(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblTransactionMgr) UpdateTransactionState(trxId int, state int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var transaction tblTransaction
	transaction.State = state
	_, err := GetDBEngine().Id(trxId).Cols("state").Update(&transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblTransactionMgr) UpdateTransactionStateFeeCost(trxId int, state int, feeCost *string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var transaction tblTransaction
	transaction.State = state
	if feeCost != nil {
		transaction.Feecost = *feeCost
	}
	_, err := GetDBEngine().Id(trxId).Update(&transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *tblTransactionMgr) GetTransactionsByState(state int) ([]tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var trxs []tblTransaction
	err := GetDBEngine().Where("state=?", state).Find(&trxs)
	return trxs, err
}

func (t *tblTransactionMgr) GetTransactions(walletId []int, coinId []int, acctId []int, state []int,
	trxTime [2]string, offSet int, limit int) (int, []tblTransaction, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	dbSession := GetDBEngine().Where("")
	if walletId != nil && len(walletId) != 0 {
		dbSession = dbSession.In("walletid", walletId)
	}
	if coinId != nil && len(coinId) != 0 {
		dbSession = dbSession.In("coinid", coinId)
	}
	if acctId != nil && len(acctId) != 0 {
		dbSession = dbSession.In("acctid", acctId)
	}
	if state != nil && len(state) != 0 {
		dbSession = dbSession.In("state", state)
	}
	if trxTime[0] != "" {
		dbSession = dbSession.And("trxtime > ?", trxTime[0])
	}
	if trxTime[1] != "" {
		dbSession = dbSession.And("trxtime < ?", trxTime[1])
	}
	var trx tblTransaction
	total, err := dbSession.Count(&trx)
	if err != nil {
		return 0, nil, err
	}

	dbSession2 := GetDBEngine().Where("")
	if walletId != nil && len(walletId) != 0 {
		dbSession2 = dbSession2.In("walletid", walletId)
	}
	if coinId != nil && len(coinId) != 0 {
		dbSession2 = dbSession2.In("coinid", coinId)
	}
	if acctId != nil && len(acctId) != 0 {
		dbSession2 = dbSession2.In("acctid", acctId)
	}
	if state != nil && len(state) != 0 {
		dbSession2 = dbSession2.In("state", state)
	}
	if trxTime[0] != "" {
		dbSession2 = dbSession2.And("trxtime > ?", trxTime[0])
	}
	if trxTime[1] != "" {
		dbSession2 = dbSession2.And("trxtime < ?", trxTime[1])
	}
	trxs := make([]tblTransaction, 0)
	dbSession2.Limit(limit, offSet).Desc("trxtime").Find(&trxs)
	return int(total), trxs, nil
}

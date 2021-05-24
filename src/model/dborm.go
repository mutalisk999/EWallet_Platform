package model

import (
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type DBMgr struct {
	DBEngine              *xorm.Engine
	SequenceMgr           *tblSequenceMgr
	PubKeyPoolMgr         *tblPubKeyPoolMgr
	CoinConfigMgr         *tblCoinConfigMgr
	AcctConfigMgr         *tblAcctConfigMgr
	WalletConfigMgr       *tblWalletConfigMgr
	AcctWalletRelationMgr *tblAcctWalletRelationMgr
	TransactionMgr        *tblTransactionMgr
	NotificationMgr       *tblNotificationMgr
	OperationLogMgr       *tblOperationLogMgr
}

var GlobalDBMgr *DBMgr

func GetDBEngine() *xorm.Engine {
	return GlobalDBMgr.DBEngine
}

func InitDB(dbType string, dbSource string) error {
	var err error
	GlobalDBMgr = new(DBMgr)
	GlobalDBMgr.DBEngine, err = xorm.NewEngine(dbType, dbSource)
	if err != nil {
		return err
	}
	GlobalDBMgr.DBEngine.SetTableMapper(core.SnakeMapper{})
	GlobalDBMgr.DBEngine.SetColumnMapper(core.SnakeMapper{})

	GlobalDBMgr.SequenceMgr = new(tblSequenceMgr)
	GlobalDBMgr.SequenceMgr.Init()

	GlobalDBMgr.PubKeyPoolMgr = new(tblPubKeyPoolMgr)
	GlobalDBMgr.PubKeyPoolMgr.Init()

	GlobalDBMgr.CoinConfigMgr = new(tblCoinConfigMgr)
	GlobalDBMgr.CoinConfigMgr.Init()

	GlobalDBMgr.AcctConfigMgr = new(tblAcctConfigMgr)
	GlobalDBMgr.AcctConfigMgr.Init()

	GlobalDBMgr.WalletConfigMgr = new(tblWalletConfigMgr)
	GlobalDBMgr.WalletConfigMgr.Init()

	GlobalDBMgr.AcctWalletRelationMgr = new(tblAcctWalletRelationMgr)
	GlobalDBMgr.AcctWalletRelationMgr.Init()

	GlobalDBMgr.TransactionMgr = new(tblTransactionMgr)
	GlobalDBMgr.TransactionMgr.Init()

	GlobalDBMgr.NotificationMgr = new(tblNotificationMgr)
	GlobalDBMgr.NotificationMgr.Init()

	GlobalDBMgr.OperationLogMgr = new(tblOperationLogMgr)
	GlobalDBMgr.OperationLogMgr.Init()

	return nil
}

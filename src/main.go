package main

import (
	"config"
	"controller"
	"fmt"
	"github.com/kataras/iris"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"model"
	"os"
	"session"
	"time"
	"utils"
	"coin"
)

var app *iris.Application
var goroutineMgr *goroutine_mgr.GoroutineManager

func DoSessionMaintain(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()
	mgr := session.GlobalSessionMgr
	sessionsOvertime := make([]string, 0)
	for {
		mgr.Mutex.Lock()
		for sid, sessionValue := range mgr.SessionStore {
			if time.Now().Unix()-sessionValue.UpdateTime.Unix() > 30*60 {
				sessionsOvertime = append(sessionsOvertime, sid)
			}
		}
		mgr.Mutex.Unlock()

		for _, sid := range sessionsOvertime {
			session.GlobalSessionMgr.DeleteSessionValue(sid)
		}
		time.Sleep(30 * time.Second)
	}
}

func DoTransactionMaintain(goroutine goroutine_mgr.Goroutine, args ...interface{}) {
	defer goroutine.OnQuit()
	for {
		trxMgr := model.GlobalDBMgr.TransactionMgr
		trxs, err := trxMgr.GetTransactionsByState(1)
		if err != nil {
			fmt.Println("DoTransactionMaintain | GetTransactionsByState: " + err.Error())
		}

		for _, trx := range trxs {
			coinConfigMgr := model.GlobalDBMgr.CoinConfigMgr
			coinConfig, err := coinConfigMgr.GetCoin(trx.Coinid)
			if err != nil {
				fmt.Println("DoTransactionMaintain | GetCoin: " + err.Error())
				continue
			}

			isConfirmed, err := coin.IsTrxConfirmed(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport,
				coinConfig.Rpcuser, coinConfig.Rpcpass, trx.Rawtrxid)
			if err != nil {
				fmt.Println("DoTransactionMaintain | IsTrxConfirmed: " + err.Error())
				continue
			}
			if isConfirmed {
				err := coin.ConfirmTransaction(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport,
					coinConfig.Rpcuser, coinConfig.Rpcpass, trx.Trxid, trx.Rawtrxid)
				if err != nil {
					fmt.Println("DoTransactionMaintain | ConfirmTransaction: " + err.Error())
					continue
				}
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(30 * time.Second)
	}
}

func StartSessionMaintainer() uint64 {
	return goroutineMgr.GoroutineCreatePn("SessionMaintainer", DoSessionMaintain, nil)
}

func StartTransactionMaintainer() uint64 {
	return goroutineMgr.GoroutineCreatePn("TransactionMaintainer", DoTransactionMaintain, nil)
}

func LoadConf() error {
	// init config
	jsonParser := new(config.JsonStruct)
	err := jsonParser.Load("config.json", &config.GlobalConfig)
	if err != nil {
		fmt.Println("Load config.json", err)
		return err
	}
	return nil
}

func Init() error {
	err := LoadConf()
	if err != nil {
		return err
	}

	utils.InitGlobalError()
	session.InitSessionMgr()

	goroutineMgr = new(goroutine_mgr.GoroutineManager)
	goroutineMgr.Initialise("MainGoroutineManager")

	fmt.Println("db path:", config.GlobalConfig.DbConfig.DbSource)
	err = model.InitDB(config.GlobalConfig.DbConfig.DbType, config.GlobalConfig.DbConfig.DbSource)
	if err != nil {
		return err
	}

	app = iris.New()
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	return nil
}

func Run(endpoint string, charset string) {
	app.Run(iris.Addr(endpoint), iris.WithCharset(charset))
}

func RunTLS(endpoint string, certFile string, keyFile string, charset string) {
	app.Run(iris.TLS(endpoint, certFile, keyFile), iris.WithCharset(charset))
}

func RegisterUrlRouter() {
	app.Post("/apis/authcode", controller.AuthCodeController)
	app.Post("/apis/identity", controller.IdentityController)
	app.Post("/apis/user", controller.UserController)
	app.Post("/apis/account", controller.AccountController)
	app.Post("/apis/wallet", controller.WalletController)
	app.Post("/apis/notification", controller.NotificationController)
	app.Post("/apis/transaction", controller.TransactionController)
	app.Post("/apis/log", controller.LogController)
	app.Post("/apis/coin", controller.CoinController)
	app.Post("/apis/manager", controller.ManagerController)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		config.IsTestEnvironment = true
		fmt.Println("Run Test Environment")
	} else {
		config.IsTestEnvironment = false
		fmt.Println("Run Product Environment")
	}
	err := Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	RegisterUrlRouter()
	StartSessionMaintainer()
	StartTransactionMaintainer()

	// http
	if !config.GlobalConfig.IsHttps {
		Run(config.GlobalConfig.HttpConfig.EndPoint, config.GlobalConfig.HttpConfig.CharSet)
	} else {
	// https
		RunTLS(config.GlobalConfig.HttpsConfig.EndPoint, config.GlobalConfig.HttpsConfig.CertFile,
			config.GlobalConfig.HttpsConfig.KeyFile, config.GlobalConfig.HttpsConfig.CharSet)
	}
}

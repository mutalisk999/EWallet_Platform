package controller

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"model"
	"session"
	"strconv"
	"strings"
	"utils"
	"coin"
	"github.com/mutalisk999/bitcoin-lib/src/utility"
	"encoding/hex"
)

const (
	LogFormatTypeTrxTransfer = 1
	LogFormatTypeTrxConfirm  = 2
	LogFormatTypeTrxRevoke   = 3
)

func CreateTransactionString(walletId int, coinId int, contractAddr string, acctId int, from string, to string,
	amount string, needConfirm int, fee string, gasPrice string, gasLimit string) (string) {
	walletIdStr := strconv.Itoa(walletId)
	coinIdStr := strconv.Itoa(coinId)
	acctIdStr := strconv.Itoa(acctId)
	needConfirmStr := strconv.Itoa(needConfirm)

	trxStr := strings.Join([]string{walletIdStr, coinIdStr, contractAddr, acctIdStr, from, to,
		amount, needConfirmStr, fee, gasPrice, gasLimit}, ",")
	return trxStr
}

func GetTransactionLogFormat(fmtType int) string {
	if fmtType == LogFormatTypeTrxTransfer {
		return "用户[%s]发起了一笔从钱包[%s]的转账,交易ID:[%s],结果:[%s]"
	} else if fmtType == LogFormatTypeTrxConfirm {
		return "用户[%s]对交易[%s]进行了确认操作,确认结果:[%s]"
	} else if fmtType == LogFormatTypeTrxRevoke {
		return "用户[%s]对交易[%s]进行了撤销操作,撤销结果:[%s]"
	}
	return ""
}

type GetWalletTrxsParam struct {
	SessionId string    `json:"sessionid"`
	WalletId  []int     `json:"walletid"`
	CoinId    []int     `json:"coinid"`
	AcctId    []int     `json:"acctid"`
	TrxTime   [2]string `json:"trxtime"`
	State     []int     `json:"state"`
	OffSet    int       `json:"offset"`
	Limit     int       `json:"limit"`
}

type GetWalletTrxsRequest struct {
	Id      int                  `json:"id"`
	JsonRpc string               `json:"jsonrpc"`
	Method  string               `json:"method"`
	Params  []GetWalletTrxsParam `json:"params"`
}

type WalletTrx struct {
	TrxId         int    `json:"trxid"`
	RawTrxId      string `json:"rawtrxid"`
	WalletId      int    `json:"walletid"`
	CoinId        int    `json:"coinid"`
	AcctId        int    `json:"acctid"`
	From          string `json:"from"`
	To            string `json:"to"`
	Amount        string `json:"amount"`
	Fee           string `json:"fee"`
	TrxTime       string `json:"trxtime"`
	NeedConfirm   int    `json:"needconfirm"`
	Confirmed     int    `json:"confirmed"`
	AcctConfirmed string `json:"acctconfirmed"`
	State         int    `json:"state"`
}

type GetWalletTrxsResult struct {
	Total int         `json:"total"`
	Trxs  []WalletTrx `json:"trxs"`
}

type GetWalletTrxsResponse struct {
	Id     int                  `json:"id"`
	Result *GetWalletTrxsResult `json:"result"`
	Error  *utils.Error         `json:"error"`
}

func GetTransactionController(ctx iris.Context, jsonRpcBody []byte) {
	var req GetWalletTrxsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res GetWalletTrxsResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	if sessionValue.Role != 0 && req.Params[0].WalletId != nil && len(req.Params[0].WalletId) != 0 {
		// 检查Wallet
		walletMgr := model.GlobalDBMgr.WalletConfigMgr
		walletConfigs, err := walletMgr.GetWalletsByIds(req.Params[0].WalletId)
		if err != nil {
			res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get wallet config")
			ctx.JSON(res)
			return
		}
		for _, walletConfig := range walletConfigs {
			if walletConfig.State != 1 {
				res.Error = utils.MakeError(200007, "wallet", walletConfig.Walletid)
				ctx.JSON(res)
				return
			}
		}
	}

	if sessionValue.Role != 0 && req.Params[0].CoinId != nil && len(req.Params[0].CoinId) != 0 {
		// 检查Coin
		coinMgr := model.GlobalDBMgr.CoinConfigMgr
		coinConfigs, err := coinMgr.GetCoins(req.Params[0].CoinId)
		if err != nil {
			res.Error = utils.MakeError(300001, coinMgr.TableName, "query", "get coin config")
			ctx.JSON(res)
			return
		}
		for _, coinConfig := range coinConfigs {
			if coinConfig.State != 1 {
				res.Error = utils.MakeError(200007, "coin", coinConfig.Coinid)
				ctx.JSON(res)
				return
			}
		}
	}

	if sessionValue.Role != 0 && req.Params[0].AcctId != nil && len(req.Params[0].AcctId) != 0 {
		// 检查Account
		acctMgr := model.GlobalDBMgr.AcctConfigMgr
		acctConfigs, err := acctMgr.GetAccountsByIds(req.Params[0].AcctId)
		if err != nil {
			res.Error = utils.MakeError(300001, acctMgr.TableName, "query", "get account config")
			ctx.JSON(res)
			return
		}
		for _, acctConfig := range acctConfigs {
			if acctConfig.State != 1 {
				res.Error = utils.MakeError(200007, "account", acctConfig.Acctid)
				ctx.JSON(res)
				return
			}
		}
	}

	walletIdsArgs := req.Params[0].WalletId
	if sessionValue.Role != 0 {
		relationMgr := model.GlobalDBMgr.AcctWalletRelationMgr
		relations, err := relationMgr.GetRelationsByAcctId(sessionValue.AcctId)
		if err != nil {
			res.Error = utils.MakeError(300001, relationMgr.TableName, "query", "get acct/wallet relation")
			ctx.JSON(res)
			return
		}
		for _, walletId := range walletIdsArgs {
			hasRelation := false
			for _, relation := range relations {
				if walletId == relation.Walletid {
					hasRelation = true
					break
				}
			}
			if hasRelation == false {
				res.Error = utils.MakeError(200008)
				ctx.JSON(res)
				return
			}
		}
		if walletIdsArgs == nil {
			walletIdsArgs = make([]int, 0)
		}
		if len(walletIdsArgs) == 0 {
			for _, relation := range relations {
				walletIdsArgs = append(walletIdsArgs, relation.Walletid)
			}
		}
		// Acct不拥有任何钱包
		if len(walletIdsArgs) == 0 {
			ctx.JSON(res)
			return
		}
	}

	trxMgr := model.GlobalDBMgr.TransactionMgr
	totalCount, transactions, err := trxMgr.GetTransactions(walletIdsArgs, req.Params[0].CoinId, req.Params[0].AcctId,
		req.Params[0].State, req.Params[0].TrxTime, req.Params[0].OffSet, req.Params[0].Limit)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		ctx.JSON(res)
		return
	}
	res.Result = new(GetWalletTrxsResult)
	res.Result.Total = totalCount
	res.Result.Trxs = make([]WalletTrx, len(transactions), len(transactions))
	for i, trx := range transactions {
		res.Result.Trxs[i] = WalletTrx{
			trx.Trxid, trx.Rawtrxid, trx.Walletid, trx.Coinid, trx.Acctid,
			trx.Fromaddr, trx.Toaddr, trx.Amount, trx.Feecost,utils.TimeToFormatString(trx.Trxtime),
			trx.Needconfirm, trx.Confirmed, trx.Acctconfirmed, trx.State}
	}
	ctx.JSON(res)
}

func CheckAccountAvailable(acctId int) *utils.Error {
	// 获取Account的配置信息
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(acctId)
	if err != nil {
		return utils.MakeError(300001, acctMgr.TableName, "query", "get account config")
	}
	if acctConfig.State != 1 {
		return utils.MakeError(200007, "account", acctConfig.Acctid)
	}
	return nil
}

func CheckRelationAvailable(acctId int, fromWalletId int) *utils.Error {
	// 获取Account关联的钱包信息
	relationMgr := model.GlobalDBMgr.AcctWalletRelationMgr
	relations, err := relationMgr.GetRelationsByAcctId(acctId)
	if err != nil {
		return utils.MakeError(300001, relationMgr.TableName, "query", "get acct/wallet relation")
	}
	hasRelation := false
	for _, relation := range relations {
		if fromWalletId == relation.Walletid {
			hasRelation = true
			break
		}
	}
	if hasRelation == false {
		return utils.MakeError(200008)
	}
	return nil
}

func CheckDestAddress(toAddr string, dstAddrConfig string, walletId int) *utils.Error {
	// 如果Destaddress非空  判断入账地址是否存在于Destaddress地址中
	if dstAddrConfig != "" || len(dstAddrConfig) != 0 {
		dstAddrList := strings.Split(dstAddrConfig, ",")
		inDstAddrs := false
		for _, dstAddr := range dstAddrList {
			if toAddr == dstAddr {
				inDstAddrs = true
				break
			}
		}
		if !inDstAddrs {
			return utils.MakeError(200010, toAddr, walletId)
		}
	}
	return nil
}

func TransferLog(isSuccQuit bool, acctId int, walletId int, trxId *int) {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(acctId)
	if err != nil {
		return
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	trxIdStr := ""
	if trxId != nil {
		trxIdStr = strconv.Itoa(*trxId)
	}

	walletMgr := model.GlobalDBMgr.WalletConfigMgr
	walletConfig, err := walletMgr.GetWalletById(walletId)
	if err != nil {
		return
	}

	logContent := fmt.Sprintf(GetTransactionLogFormat(LogFormatTypeTrxTransfer), acctConfig.Realname,
		walletConfig.Walletname, trxIdStr, resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(acctId, 5, logContent)
	if err != nil {
		return
	}
	return
}

type CreateTrxParam struct {
	SessionId    string `json:"sessionid"`
	OperateId    int    `json:"operateid"`
	FromWalletId int    `json:"fromwalletid"`
	ToAddr       string `json:"toaddr"`
	Amount       string `json:"amount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	Signature    string `json:"signature"`
}

type CreateTrxRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []CreateTrxParam `json:"params"`
}

type CreateTrxResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func CreateTrxController(ctx iris.Context, jsonRpcBody []byte) {
	var req CreateTrxRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res CreateTrxResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	if sessionValue.Role != 1 {
		res.Error = utils.MakeError(200009)
		ctx.JSON(res)
		return
	}

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(4, req.Params[0].OperateId)
	if !verify || err != nil {
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	// 验证签名
	funcNameStr := "transfer"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	fromWalletIdStr := strconv.Itoa(req.Params[0].FromWalletId)
	toAddrStr := req.Params[0].ToAddr
	amountStr := req.Params[0].Amount
	feeStr := req.Params[0].Fee
	gasPriceStr := req.Params[0].GasPrice
	gasLimitStr := req.Params[0].GasLimit
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, fromWalletIdStr, toAddrStr, amountStr, feeStr, gasPriceStr, gasLimitStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	// 检查Account的有效性
	uErr := CheckAccountAvailable(sessionValue.AcctId)
	if uErr != nil {
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 检查交易出账的Wallet是否属于Account关联的Wallet
	uErr = CheckRelationAvailable(sessionValue.AcctId, req.Params[0].FromWalletId)
	if uErr != nil {
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 获取Wallet的配置信息
	walletMgr := model.GlobalDBMgr.WalletConfigMgr
	walletConfig, err := walletMgr.GetWalletById(req.Params[0].FromWalletId)
	if err != nil {
		res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get wallet config")
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	if walletConfig.State != 1 {
		res.Error = utils.MakeError(200007, "wallet", walletConfig.Walletid)
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 获取钱包对应PubKey信息
	keyPoolMgr := model.GlobalDBMgr.PubKeyPoolMgr
	pubkeyStr, err := keyPoolMgr.QueryPubKeyByKeyIndex(walletConfig.Keyindex)
	if err != nil {
		res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get key pool config")
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 获取Coin配置信息
	coinMgr := model.GlobalDBMgr.CoinConfigMgr
	coinConfig, err := coinMgr.GetCoin(walletConfig.Coinid)
	if err != nil {
		res.Error = utils.MakeError(300001, coinMgr.TableName, "query", "get coin config")
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	if coinConfig.State != 1 {
		res.Error = utils.MakeError(200007, "coin", coinConfig.Coinid)
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	if !config.IsSupportCoin(coinConfig.Coinsymbol) {
		res.Error = utils.MakeError(600001, coinConfig.Coinsymbol)
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	coinConfigDetail, _ := config.GlobalSupportCoinMgr[coinConfig.Coinsymbol]

	// 检查to address的合法性
	addrValid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser,
		coinConfig.Rpcpass, req.Params[0].ToAddr)
	if err != nil {
		res.Error = utils.MakeError(800000, err.Error())
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	if !addrValid {
		res.Error = utils.MakeError(500003, req.Params[0].ToAddr)
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 检查存在Dest Address设置下  交易中to address的有效性
	uErr = CheckDestAddress(req.Params[0].ToAddr, walletConfig.Destaddress, walletConfig.Walletid)
	if uErr != nil {
		res.Error = uErr
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 生成交易防篡改签名
	trxStr := CreateTransactionString(req.Params[0].FromWalletId, walletConfig.Coinid, coinConfigDetail.ContractAddress,
		sessionValue.AcctId, walletConfig.Address, req.Params[0].ToAddr, req.Params[0].Amount, walletConfig.Needsigcount,
		req.Params[0].Fee, req.Params[0].GasPrice, req.Params[0].GasLimit)
	signature, err := coin.CoinSignTrx('2', utility.Sha256([]byte(trxStr)), uint16(walletConfig.Keyindex))
	if err != nil {
		res.Error = utils.MakeError(900002)
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	signatureHex := hex.EncodeToString(signature)

	trxMgr := model.GlobalDBMgr.TransactionMgr
	trxId, err := trxMgr.NewTransaction(req.Params[0].FromWalletId, walletConfig.Coinid, coinConfigDetail.ContractAddress,
		sessionValue.AcctId, walletConfig.Address, req.Params[0].ToAddr, req.Params[0].Amount, walletConfig.Needsigcount,
		req.Params[0].Fee, req.Params[0].GasPrice, req.Params[0].GasLimit, signatureHex)

	trx, err := trxMgr.GetTransactionById(trxId)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	if walletConfig.Needsigcount == 1 {
		feeCostStr, trxId, err := coin.Transfer(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass,
			uint16(walletConfig.Keyindex), pubkeyStr, walletConfig.Address, req.Params[0].ToAddr, req.Params[0].Amount,
			req.Params[0].Fee, req.Params[0].GasPrice, req.Params[0].GasLimit)
		if err != nil {
			res.Error = utils.MakeError(800000, err.Error())
			TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
			ctx.JSON(res)
			return
		}

		if feeCostStr != ""{
			trx.Feecost = feeCostStr
		}
		trx.Rawtrxid = trxId
		trx.State = 1
	}

	trx.Confirmed = 1
	trx.Acctconfirmed = strconv.Itoa(sessionValue.AcctId)

	err = trxMgr.UpdateTransaction(trx)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "update", "update transaction")
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}

	// 创建提醒
	relationMgr := model.GlobalDBMgr.AcctWalletRelationMgr
	relations, err := relationMgr.GetRelationsByWalletId(req.Params[0].FromWalletId)
	if err != nil {
		res.Error = utils.MakeError(300001, relationMgr.TableName, "query", "get acct/wallet relation")
		TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
		ctx.JSON(res)
		return
	}
	notifyMgr := model.GlobalDBMgr.NotificationMgr
	for _, relation := range relations {
		if sessionValue.AcctId != relation.Acctid {
			_, err := notifyMgr.NewNotification(&relation.Acctid, &req.Params[0].FromWalletId, &trx.Trxid,
				1, fmt.Sprintf("有一笔新的转账交易产生, 需要您去处理, 交易ID: %d", trx.Trxid),
				0, "", "")
			if err != nil {
				res.Error = utils.MakeError(300001, trxMgr.TableName, "insert", "insert notification")
				TransferLog(false, sessionValue.AcctId, req.Params[0].FromWalletId, nil)
				ctx.JSON(res)
				return
			}
		}
	}

	TransferLog(true, sessionValue.AcctId, req.Params[0].FromWalletId, &trx.Trxid)
	ctx.JSON(res)
	return
}

func RevokeLog(isSuccQuit bool, acctId int, trxId int) {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(acctId)
	if err != nil {
		return
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	logContent := fmt.Sprintf(GetTransactionLogFormat(LogFormatTypeTrxRevoke), acctConfig.Realname,
		strconv.Itoa(trxId), resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(acctId, 5, logContent)
	if err != nil {
		return
	}
	return
}

type RevokeTrxParam struct {
	SessionId string `json:"sessionid"`
	OperateId int    `json:"operateid"`
	TrxId     int    `json:"trxid"`
	Signature string `json:"signature"`
}

type RevokeTrxRequest struct {
	Id      int              `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  []RevokeTrxParam `json:"params"`
}

type RevokeTrxResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func RevokeTrxController(ctx iris.Context, jsonRpcBody []byte) {
	var req RevokeTrxRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res RevokeTrxResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(4, req.Params[0].OperateId)
	if !verify || err != nil {
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	// 验证签名
	funcNameStr := "revoke"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	trxIdStr := strconv.Itoa(req.Params[0].TrxId)
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, trxIdStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	trxMgr := model.GlobalDBMgr.TransactionMgr
	trx, err := trxMgr.GetTransactionById(req.Params[0].TrxId)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		ctx.JSON(res)
		return
	}
	if trx.State != 0 {
		res.Error = utils.MakeError(200007, "trx", trx.Trxid)
		ctx.JSON(res)
		return
	}

	if sessionValue.Role == 1 {
		if sessionValue.AcctId != trx.Acctid {
			res.Error = utils.MakeError(200011)
			ctx.JSON(res)
			return
		}
	}

	err = trxMgr.UpdateTransactionState(req.Params[0].TrxId, 3)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "update", "update transaction state")
		RevokeLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 删除提醒
	notifyMgr := model.GlobalDBMgr.NotificationMgr
	notifyMgr.DeleteNotification(nil, nil, nil, &trx.Trxid, nil, nil, nil, nil)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "delete", "delete notification")
		RevokeLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	RevokeLog(true, sessionValue.AcctId, req.Params[0].TrxId)
	ctx.JSON(res)
	return
}

func ConfirmLog(isSuccQuit bool, acctId int, trxId int) {
	acctMgr := model.GlobalDBMgr.AcctConfigMgr
	acctConfig, err := acctMgr.GetAccountById(acctId)
	if err != nil {
		return
	}

	resultStr := "失败"
	if isSuccQuit {
		resultStr = "成功"
	}

	logContent := fmt.Sprintf(GetTransactionLogFormat(LogFormatTypeTrxConfirm), acctConfig.Realname,
		strconv.Itoa(trxId), resultStr)
	logMgr := model.GlobalDBMgr.OperationLogMgr
	_, err = logMgr.NewOperatorLog(acctId, 5, logContent)
	if err != nil {
		return
	}
	return
}

type ConfirmTrxParam struct {
	SessionId string `json:"sessionid"`
	OperateId int    `json:"operateid"`
	TrxId     int    `json:"trxid"`
	Signature string `json:"signature"`
}

type ConfirmTrxRequest struct {
	Id      int               `json:"id"`
	JsonRpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []ConfirmTrxParam `json:"params"`
}

type ConfirmTrxResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}

func ConfirmTrxController(ctx iris.Context, jsonRpcBody []byte) {
	var req ConfirmTrxRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res ConfirmTrxResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	sessionValue, ok := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !ok {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}
	if sessionValue.Role != 1 {
		res.Error = utils.MakeError(200012)
		ctx.JSON(res)
		return
	}

	verify, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(4, req.Params[0].OperateId)
	if !verify || err != nil {
		res.Error = utils.MakeError(400005)
		ctx.JSON(res)
		return
	}

	// 验证签名
	funcNameStr := "confirm"
	sessionIdStr := req.Params[0].SessionId
	operatorIdStr := strconv.Itoa(req.Params[0].OperateId)
	trxIdStr := strconv.Itoa(req.Params[0].TrxId)
	sigSrcStr := strings.Join([]string{funcNameStr, sessionIdStr, operatorIdStr, trxIdStr}, ",")
	err = utils.RsaVerySignWithSha1Hex(sigSrcStr, req.Params[0].Signature, sessionValue.PubKey)
	if err != nil {
		res.Error = utils.MakeError(400002)
		ctx.JSON(res)
		return
	}

	trxMgr := model.GlobalDBMgr.TransactionMgr
	trx, err := trxMgr.GetTransactionById(req.Params[0].TrxId)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "query", "get transaction")
		ctx.JSON(res)
		return
	}
	if trx.State != 0 {
		res.Error = utils.MakeError(200007, "trx", trx.Trxid)
		ctx.JSON(res)
		return
	}

	// 达到最大确认数
	if trx.Confirmed == trx.Needconfirm {
		res.Error = utils.MakeError(200013)
		ctx.JSON(res)
		return
	}

	acctConfirmed := strings.Split(trx.Acctconfirmed, ",")
	for _, acctStr := range acctConfirmed {
		if strconv.Itoa(sessionValue.AcctId) == acctStr {
			res.Error = utils.MakeError(200014)
			ctx.JSON(res)
			return
		}
	}

	// 检查Account的有效性
	uErr := CheckAccountAvailable(sessionValue.AcctId)
	if uErr != nil {
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 检查交易出账的Wallet是否属于Account关联的Wallet
	uErr = CheckRelationAvailable(sessionValue.AcctId, trx.Walletid)
	if uErr != nil {
		res.Error = uErr
		ctx.JSON(res)
		return
	}

	// 获取Wallet的配置信息
	walletMgr := model.GlobalDBMgr.WalletConfigMgr
	walletConfig, err := walletMgr.GetWalletById(trx.Walletid)
	if err != nil {
		res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get wallet config")
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}
	if walletConfig.State != 1 {
		res.Error = utils.MakeError(200007, "wallet", walletConfig.Walletid)
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 获取钱包对应PubKey信息
	keyPoolMgr := model.GlobalDBMgr.PubKeyPoolMgr
	pubkeyStr, err := keyPoolMgr.QueryPubKeyByKeyIndex(walletConfig.Keyindex)
	if err != nil {
		res.Error = utils.MakeError(300001, walletMgr.TableName, "query", "get key pool config")
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 获取Coin配置信息
	coinMgr := model.GlobalDBMgr.CoinConfigMgr
	coinConfig, err := coinMgr.GetCoin(walletConfig.Coinid)
	if err != nil {
		res.Error = utils.MakeError(300001, coinMgr.TableName, "query", "get coin config")
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}
	if coinConfig.State != 1 {
		res.Error = utils.MakeError(200007, "coin", coinConfig.Coinid)
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	if !config.IsSupportCoin(coinConfig.Coinsymbol) {
		res.Error = utils.MakeError(600001, coinConfig.Coinsymbol)
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 检查存在Dest Address设置下  交易中to address的有效性
	uErr = CheckDestAddress(trx.Toaddr, walletConfig.Destaddress, walletConfig.Walletid)
	if uErr != nil {
		res.Error = uErr
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	if trx.Confirmed + 1 == trx.Needconfirm {
		feeCostStr, trxId, err := coin.Transfer(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass,
			uint16(walletConfig.Keyindex), pubkeyStr, walletConfig.Address, trx.Toaddr, trx.Amount, trx.Fee, trx.Gasprice, trx.Gaslimit)
		if err != nil {
			res.Error = utils.MakeError(800000, err.Error())
			ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
			ctx.JSON(res)
			return
		}

		if feeCostStr != ""{
			trx.Feecost = feeCostStr
		}
		trx.Rawtrxid = trxId
		trx.State = 1
	}

	// 交易防篡改签名验签
	trxStr := CreateTransactionString(trx.Walletid, trx.Coinid, trx.Contractaddr,
		trx.Acctid, trx.Fromaddr, trx.Toaddr, trx.Amount, trx.Needconfirm,
		trx.Fee, trx.Gasprice, trx.Gaslimit)
	signature, err := hex.DecodeString(trx.Signature)
	if err != nil {
		res.Error = utils.MakeError(900004)
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}
	_, err = coin.CoinVerifyTrx('2', uint16(walletConfig.Keyindex), utility.Sha256([]byte(trxStr)), signature)
	if err != nil {
		res.Error = utils.MakeError(900003)
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	trx.Trxid = req.Params[0].TrxId
	trx.Confirmed = trx.Confirmed + 1
	trx.Acctconfirmed = trx.Acctconfirmed + "," + strconv.Itoa(sessionValue.AcctId)

	err = trxMgr.UpdateTransaction(trx)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "update", "update transaction")
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	// 删除提醒
	notifyMgr := model.GlobalDBMgr.NotificationMgr
	var pAcctId *int
	pAcctId = nil
	if trx.State == 0 {
		pAcctId = &sessionValue.AcctId
	}
	notifyMgr.DeleteNotification(nil, pAcctId, nil, &trx.Trxid, nil, nil, nil, nil)
	if err != nil {
		res.Error = utils.MakeError(300001, trxMgr.TableName, "delete", "delete notification")
		ConfirmLog(false, sessionValue.AcctId, req.Params[0].TrxId)
		ctx.JSON(res)
		return
	}

	ConfirmLog(true, sessionValue.AcctId, req.Params[0].TrxId)
	ctx.JSON(res)
	return

}

func TransactionController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}

	var res utils.JsonRpcResponse
	if funcName == "get_wallet_trxs" {
		GetTransactionController(ctx, jsonRpcBody)
	} else if funcName == "transfer" {
		CreateTrxController(ctx, jsonRpcBody)
	} else if funcName == "confirm" {
		ConfirmTrxController(ctx, jsonRpcBody)
	} else if funcName == "revoke" {
		RevokeTrxController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}
}

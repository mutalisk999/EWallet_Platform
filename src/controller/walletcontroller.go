package controller

import (
	"coin"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"model"
	"session"
	"strconv"
	"utils"
	"strings"
)

type EmptyResponse struct {
	Id     int          `json:"id"`
	Result *int         `json:"result"`
	Error  *utils.Error `json:"error"`
}
type RequestBase struct {
	Id      int    `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
}
type ListWalletsParam struct {
	SessionId string `json:"sessionid"`
	AcctIds   []int  `json:"acctids"`
	CoinId    []int  `json:"coinid"`
	State     []int  `json:"state"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}
type ListWalletsRequest struct {
	RequestBase
	Params []ListWalletsParam `json:"params"`
}
type WalletResult struct {
	WalletId   int    `json:"walletid"`
	WalletName string `json:"walletname"`
	CoinId     int    `json:"coinid"`
	CoinSymbol string `json:"coinsymbol"`
	Balance    string `json:"balance"`
	FeeBalance string `json:"feebalance"`
	Address    string `json:"address"`
	State      int    `json:"state"`
}
type ListWalletResult struct {
	Total   int64          `json:"total"`
	Wallets []WalletResult `json:"wallets"`
}
type ListWalletResponse struct {
	Id     int               `json:"id"`
	Result *ListWalletResult `json:"result"`
	Error  *utils.Error      `json:"error"`
}

func BasicCheck(ctx iris.Context, sessionid string, role []int, seqtype int, mgmtid int, signature string, orgindata string) (bool, int, *utils.Error) {
	//SessionCheck
	sessionval, exist := session.GlobalSessionMgr.GetSessionValue(sessionid)
	if !exist {
		return false, 0, utils.MakeError(200004)
	}
	//defer session.GlobalSessionMgr.RefreshSessionValue(sessionid)
	//role Check
	if len(role) != 0 {
		inrole := false
		for erole := range role {
			if erole == sessionval.Role {
				inrole = true
				break
			}
		}
		if !inrole {
			return false, sessionval.AcctId, utils.MakeError(400009)
		}
	}
	//mgmtid Check
	if seqtype != 0 {
		vres, err := model.GlobalDBMgr.SequenceMgr.VerifySequence(seqtype, mgmtid)
		if !vres || err != nil {
			return false, sessionval.AcctId, utils.MakeError(400010)
		}
	}

	//Signature Check
	if orgindata != "" {
		err := utils.RsaVerySignWithSha1Hex(orgindata, signature, sessionval.PubKey)
		if err != nil {
			return false, sessionval.AcctId, utils.MakeError(400002)
		}
	}
	return true, sessionval.AcctId, nil

}
func ListWalletController(ctx iris.Context, jsonRpcBody []byte) {
	var req ListWalletsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res ListWalletResponse
	res.Id = req.Id
	res.Result = nil
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}

	seval, exist := session.GlobalSessionMgr.GetSessionValue(req.Params[0].SessionId)
	if !exist {
		res.Error = utils.MakeError(200004)
		ctx.JSON(res)
		return
	}

	cids := make([]int, 0)
	for _, coinid := range req.Params[0].CoinId {
		_, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(coinid)
		if err != nil {
			if err.Error() != "key not found" {
				res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "query", "list wallets")
				ctx.JSON(res)
				return
			} else {
				res.Error = utils.MakeError(500001, coinid)
				ctx.JSON(res)
				return
			}
		}
		cids = append(cids, coinid)
	}
	var accids = make([]int, 0)
	if seval.Role == 0 {
		accids = append(accids, req.Params[0].AcctIds...)
	} else {
		accids = append(accids, seval.AcctId)
	}
	wallets, total, werr := model.GlobalDBMgr.WalletConfigMgr.ListWallets(cids, req.Params[0].State, accids, req.Params[0].Offset, req.Params[0].Limit)
	if werr != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "query", "list wallets")
		ctx.JSON(res)
		return
	}
	res.Result = new(ListWalletResult)
	res.Result.Total = total
	for _, wa := range wallets {
		var walletres WalletResult
		walletres.State = wa.State
		walletres.Address = wa.Address
		walletres.CoinId = wa.Coinid
		coincfg, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(wa.Coinid)
		if coincfg.State != 1 || err != nil {
			continue
		}

		ba, fee_balance, err := coin.GetBalance(coincfg.Coinsymbol, coincfg.Ip, coincfg.Rpcport, coincfg.Rpcuser, coincfg.Rpcpass, wa.Address)
		if err != nil {
			fmt.Println("no money")
			//continue
		}
		walletres.Balance = ba
		walletres.FeeBalance = fee_balance
		walletres.CoinSymbol = coincfg.Coinsymbol
		walletres.WalletId = wa.Walletid
		walletres.WalletName = wa.Walletname
		res.Result.Wallets = append(res.Result.Wallets, walletres)
	}
	ctx.JSON(res)
	return

}

type CreateWalletParam struct {
	SessionId    string `json:"sessionid"`
	MgmtId       int    `json:"mgmtid"`
	WalletName   string `json:"walletname"`
	CoinId       int    `json:"coinid"`
	DestAddress  string `json:"destaddress"`
	NeedSigCount int    `json:"needsigcount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	SigUserId    []int  `json:"siguserid"`
	State        int    `json:"state"`
	Signature    string `json:"signature"`
}
type CreateWalletRequest struct {
	RequestBase
	Params []CreateWalletParam `json:"params"`
}

func CreateWalletController(ctx iris.Context, jsonRpcBody []byte) {
	var req CreateWalletRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	if pa.NeedSigCount <= 0 {
		res.Error = utils.MakeError(500004)
		ctx.JSON(res)
		return
	}
	orgindata := "create_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + pa.WalletName + "," + strconv.Itoa(pa.CoinId) +
		"," + pa.DestAddress + "," + strconv.Itoa(pa.NeedSigCount) + "," + pa.Fee + "," + pa.GasPrice + "," + pa.GasLimit + "," +
		utils.IntArrayToString(pa.SigUserId) + "," + strconv.Itoa(pa.State)
	checkres, acctid, errres := BasicCheck(ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		res.Error = errres
		ctx.JSON(res)
		return
	}
	keyindex, err := model.GlobalDBMgr.PubKeyPoolMgr.GetAnUnusedKeyIndex()
	if err != nil {
		if err.Error() == "query key error!" {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.PubKeyPoolMgr.TableName, "query", "query key!")
			ctx.JSON(res)
			return
		} else if err.Error() == "key not found!" {
			res.Error = utils.MakeError(700000)
			ctx.JSON(res)
			return
		} else {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.PubKeyPoolMgr.TableName, "update", "set key used")
			ctx.JSON(res)
			return
		}
	}
	pubKeyStr, err := model.GlobalDBMgr.PubKeyPoolMgr.UsePubkey(keyindex)
	if err != nil {
		if err.Error() == "key not found!" {
			res.Error = utils.MakeError(700000)
			ctx.JSON(res)
			return
		} else {
			res.Error = utils.MakeError(700000)
			ctx.JSON(res)
			return
		}
	}
	coinConfig, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(pa.CoinId)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get coin by id")
		ctx.JSON(res)
		return
	}
	address, err := coin.GetAddressByPubKey(pubKeyStr, coinConfig.Coinsymbol)
	if err != nil {
		res.Error = utils.MakeError(200016)
		ctx.JSON(res)
		return
	}

	if pa.DestAddress != "" {
		dstAddrList := strings.Split(pa.DestAddress, ",")
		for _, dstAddress := range dstAddrList {
			valid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, dstAddress)
			if err != nil {
				res.Error = utils.MakeError(800000, err.Error())
				ctx.JSON(res)
				return
			}
			if !valid {
				res.Error = utils.MakeError(500003, dstAddress)
				ctx.JSON(res)
				return
			}
		}
	}

	err = coin.ImportAddress(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, address)
	if err != nil {
		res.Error = utils.MakeError(800000, err.Error())
		ctx.JSON(res)
		return
	}

	err = model.GlobalDBMgr.WalletConfigMgr.InsertWallet(pa.CoinId, pa.WalletName, keyindex, address, pa.DestAddress,
		pa.NeedSigCount, pa.Fee, pa.GasPrice, pa.GasLimit, pa.State)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Insert", "create wallets")
		ctx.JSON(res)
		return
	}
	wa, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletByName(pa.WalletName)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get wallet by name")
		ctx.JSON(res)
		return
	}
	for _,uid := range pa.SigUserId {
		err = model.GlobalDBMgr.AcctWalletRelationMgr.InsertRelation(uid, wa.Walletid)
		if err != nil {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "Insert", "InsertRelation")
			ctx.JSON(res)
			return
		}
	}
	logmsg := "创建钱包,钱包名:" + pa.WalletName + ",CoinId:" + strconv.Itoa(pa.CoinId) + ",DestAddress:" + pa.DestAddress +
		",NeedSigCount:" + strconv.Itoa(pa.NeedSigCount) + ",Fee:" + pa.Fee + ",GasPrice:" + pa.GasPrice + ",GasLimit:" + pa.GasLimit +
		",SigUserId:" + utils.IntArrayToString(pa.SigUserId) + ",State:" + strconv.Itoa(pa.State)

	model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(acctid, 4, logmsg)

	ctx.JSON(res)
	return
}

type GetWalletParam struct {
	SessionId string `json:"sessionid"`
	WalletId  int    `json:"walletid"`
}
type GetWalletsRequest struct {
	RequestBase
	Params []GetWalletParam `json:"params"`
}
type GetWalletResult struct {
	WalletId     int    `json:"walletid"`
	WalletName   string `json:"walletname"`
	CoinId       int    `json:"coinid"`
	Address      string `json:"address"`
	DestAddress  string `json:"destaddress"`
	NeedSigCount int    `json:"needsigcount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	SigUserId    []int  `json:"siguserid"`
	State        int    `json:"state"`
}
type GetWalletResponse struct {
	Id     int             `json:"id"`
	Result GetWalletResult `json:"result"`
	Error  *utils.Error    `json:"error"`
}

func GetWalletController(ctx iris.Context, jsonRpcBody []byte) {
	var req GetWalletsRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	var res GetWalletResponse
	res.Id = req.Id
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		ctx.JSON(res)
		return
	}

	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	checkres, _, errres := BasicCheck(ctx, pa.SessionId, []int{0, 1}, 0, 0, "", "")
	if !checkres {
		res.Error = errres
		ctx.JSON(res)
		return
	}
	wa, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletById(pa.WalletId)
	if err != nil {
		if err.Error() == "no find wallet" {
			res.Error = utils.MakeError(500000, pa.WalletId)
		} else {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "query", "get wallet by id")

		}
	} else {
		res.Result.WalletId = wa.Walletid
		res.Result.WalletName = wa.Walletname
		res.Result.CoinId = wa.Coinid
		res.Result.Address = wa.Address
		res.Result.DestAddress = wa.Destaddress
		res.Result.NeedSigCount = wa.Needsigcount
		res.Result.Fee = wa.Fee
		res.Result.GasPrice = wa.Gasprice
		res.Result.GasLimit = wa.Gaslimit
		res.Result.State = wa.State
		uids := make([]int, 0)
		relations, err := model.GlobalDBMgr.AcctWalletRelationMgr.GetRelationsByWalletId(wa.Walletid)
		if err != nil {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "query", "get relation by wallet id")
			ctx.JSON(res)
			return
		}
		for _, rel := range relations {
			uids = append(uids, rel.Acctid)
		}
		res.Result.SigUserId = uids
	}
	ctx.JSON(res)
	return
}

type ModifyWalletParam struct {
	SessionId    string `json:"sessionid"`
	MgmtId       int    `json:"mgmtid"`
	Walletid     int    `json:"walletid"`
	WalletName   string `json:"walletname"`
	DestAddress  string `json:"destaddress"`
	NeedSigCount int    `json:"needsigcount"`
	Fee          string `json:"fee"`
	GasPrice     string `json:"gasprice"`
	GasLimit     string `json:"gaslimit"`
	SigUserId    []int  `json:"siguserid"`
	State        int    `json:"state"`
	Signature    string `json:"signature"`
}
type ModifyWalletRequest struct {
	RequestBase
	Params []ModifyWalletParam `json:"params"`
}

func ModifyWalletController(ctx iris.Context, jsonRpcBody []byte) {
	var req ModifyWalletRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	if pa.NeedSigCount <= 0 {
		res.Error = utils.MakeError(500004)
		ctx.JSON(res)
		return
	}
	if pa.State != 0 && pa.State != 1 && pa.State != 2 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	orgindata := "modify_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid) + "," + pa.WalletName + "," +
		pa.DestAddress + "," + strconv.Itoa(pa.NeedSigCount) + "," + pa.Fee + "," + pa.GasPrice + "," + pa.GasLimit + "," +
		utils.IntArrayToString(pa.SigUserId) + "," + strconv.Itoa(pa.State)

	checkres, acctid, errres := BasicCheck(ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		res.Error = errres
		ctx.JSON(res)
		return
	}
	wa, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletById(pa.Walletid)
	if err != nil {
		res.Error = utils.MakeError(500000, pa.Walletid)
		ctx.JSON(res)
		return
	}
	if wa.State == 3 {
		res.Error = utils.MakeError(500005, pa.Walletid)
		ctx.JSON(res)
		return
	}
	coinConfig, err := model.GlobalDBMgr.CoinConfigMgr.GetCoin(wa.Coinid)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "Query", "get coin by id")
		ctx.JSON(res)
		return
	}

	if pa.DestAddress != "" {
		dstAddrList := strings.Split(pa.DestAddress, ",")
		for _, dstAddress := range dstAddrList {
			valid, err := coin.IsAddressValid(coinConfig.Coinsymbol, coinConfig.Ip, coinConfig.Rpcport, coinConfig.Rpcuser, coinConfig.Rpcpass, dstAddress)
			if err != nil {
				res.Error = utils.MakeError(800000, err.Error())
				ctx.JSON(res)
				return
			}
			if !valid {
				res.Error = utils.MakeError(500003, dstAddress)
				ctx.JSON(res)
				return
			}
		}
	}

	err = model.GlobalDBMgr.WalletConfigMgr.UpdateWallet(pa.Walletid, pa.WalletName, pa.DestAddress, pa.NeedSigCount,
		pa.Fee, pa.GasPrice, pa.GasLimit, pa.State)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "update", "update wallet")
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.AcctWalletRelationMgr.DeleteRelationByWalletId(pa.Walletid)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "Delete", "Delete Relation by wallet id")
		ctx.JSON(res)
		return
	}
	for _, uid := range pa.SigUserId {
		err = model.GlobalDBMgr.AcctWalletRelationMgr.InsertRelation(uid, pa.Walletid)
		if err != nil {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "Insert", "InsertRelation")
			ctx.JSON(res)
			return
		}
	}
	ctx.JSON(res)
	logmsg := "修改钱包,钱包id:" + strconv.Itoa(pa.Walletid) + ",修改后属性 钱包名:" + pa.WalletName + ",DestAddress:" + pa.DestAddress +
		",NeedSigCount:" + strconv.Itoa(pa.NeedSigCount) + ",Fee:" + pa.Fee + ",GasPrice:" + pa.GasPrice + ",GasLimit:" + pa.GasLimit +
		",SigUserId:" + utils.IntArrayToString(pa.SigUserId) + ",State:" + strconv.Itoa(pa.State)

	model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(acctid, 4, logmsg)
	return
}

type DeleteWalletParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	Walletid  int    `json:"walletid"`
	Signature string `json:"signature"`
}
type DeleteWalletRequest struct {
	RequestBase
	Params []DeleteWalletParam `json:"params"`
}

func DeleteWalletController(ctx iris.Context, jsonRpcBody []byte) {
	var req DeleteWalletRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	orgindata := "delete_wallet," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid)
	checkres, acctid, errres := BasicCheck(ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		res.Error = errres
		ctx.JSON(res)
		return
	}
	wal, err := model.GlobalDBMgr.WalletConfigMgr.GetWalletById(pa.Walletid)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "select", "select wallet by id")
		ctx.JSON(res)
		return
	}
	wal.State = 3
	err = model.GlobalDBMgr.WalletConfigMgr.UpdateWallet(wal.Walletid, wal.Walletname, wal.Destaddress, wal.Needsigcount, wal.Fee, wal.Gasprice, wal.Gaslimit, wal.State)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "update", "set wallet to delete")
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.AcctWalletRelationMgr.DeleteRelationByWalletId(pa.Walletid)
	if err != nil {
		res.Error = utils.MakeError(300001, model.GlobalDBMgr.AcctWalletRelationMgr.TableName, "delete", "delete relation")

	}
	ctx.JSON(res)
	logmsg := "删除钱包,钱包id:" + strconv.Itoa(pa.Walletid)

	model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(acctid, 4, logmsg)
	return
}

type ChangeWalletStateParam struct {
	SessionId string `json:"sessionid"`
	MgmtId    int    `json:"mgmtid"`
	Walletid  int    `json:"walletid"`
	State     int    `json:"state"`
	Signature string `json:"signature"`
}
type ChangeWalletStateRequest struct {
	RequestBase
	Params []ChangeWalletStateParam `json:"params"`
}

func ChangeWalletStateController(ctx iris.Context, jsonRpcBody []byte) {
	var req ChangeWalletStateRequest
	err := json.Unmarshal(jsonRpcBody, &req)
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		return
	}
	var res EmptyResponse
	res.Id = req.Id
	if len(req.Params) != 1 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	pa := req.Params[0]
	if pa.State != 0 && pa.State != 1 && pa.State != 2 {
		res.Error = utils.MakeError(200001)
		ctx.JSON(res)
		return
	}
	orgindata := "change_wallet_state," + pa.SessionId + "," + strconv.Itoa(pa.MgmtId) + "," + strconv.Itoa(pa.Walletid) + "," + strconv.Itoa(pa.State)
	checkres, acctid, errres := BasicCheck(ctx, pa.SessionId, []int{0}, 3, pa.MgmtId, pa.Signature, orgindata)
	if !checkres {
		res.Error = errres
		ctx.JSON(res)
		return
	}
	err = model.GlobalDBMgr.WalletConfigMgr.ChangeWalletState(pa.Walletid, pa.State)
	if err != nil {
		if err.Error() == "no find wallet" {
			res.Error = utils.MakeError(500000, pa.Walletid)
		} else {
			res.Error = utils.MakeError(300001, model.GlobalDBMgr.WalletConfigMgr.TableName, "update", "change wallet state")
		}
	}
	ctx.JSON(res)
	logmsg := "更改钱包状态,钱包id:" + strconv.Itoa(pa.Walletid) + ",新状态:" + strconv.Itoa(pa.State)

	model.GlobalDBMgr.OperationLogMgr.NewOperatorLog(acctid, 4, logmsg)
	return
}
func WalletController(ctx iris.Context) {
	id, funcName, jsonRpcBody, err := utils.ReadJsonRpcBody(ctx)
	var res utils.JsonRpcResponse
	if err != nil {
		utils.SetInternalError(ctx, err.Error())
		ctx.JSON(res)
		return
	}

	if funcName == "list_wallets" {
		ListWalletController(ctx, jsonRpcBody)
	} else if funcName == "create_wallet" {
		CreateWalletController(ctx, jsonRpcBody)
	} else if funcName == "get_wallet" {
		GetWalletController(ctx, jsonRpcBody)
	} else if funcName == "modify_wallet" {
		ModifyWalletController(ctx, jsonRpcBody)
	} else if funcName == "delete_wallet" {
		DeleteWalletController(ctx, jsonRpcBody)
	} else if funcName == "change_wallet_state" {
		ChangeWalletStateController(ctx, jsonRpcBody)
	} else {
		res.Id = id
		res.Result = nil
		res.Error = utils.MakeError(200000, funcName, ctx.Path())
		ctx.JSON(res)
	}

}

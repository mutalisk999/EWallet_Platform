package coin

import (
	"github.com/mutalisk999/go-lib/src/net/buffer_tcp"
	"github.com/mutalisk999/go-cryptocard/src/cryptocard"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"config"
	"model"
	"utils"
)


func GetAddressByPubKey(pubKeyStr string, coinSymbol string) (string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "", utils.InvalidCoinSymbol(coinSymbol)
	}

	if coinSymbol == "BTC" {
		addrStr, err := BTCCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	}else if coinSymbol == "ETH"{
		addrStr,err := ETHCalcAddressByPubkey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinSymbol == "LTC" {
		addrStr, err := LTCCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinSymbol == "BCH" {
		addrStr, err := BCHCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinObj.IsErc20{
		addrStr, err := ETHCalcAddressByPubkey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinObj.IsOmni {
		addrStr, err := OMNICalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	} else if coinObj.CoinSymbol == "UB" {
		addrStr, err := UBCalcAddressByPubKey(pubKeyStr)
		if err != nil {
			return "", err
		}
		return addrStr, nil
	}
	return "", nil
}

func ImportAddress(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string) (error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return  utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		err := agent.ImportAddressRPC(address)
		if err != nil {
			return err
		}
	} else if coinSymbol == "ETH"||coinObj.IsErc20{
		return nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		err := agent.ImportAddressRPC(address)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsTrxConfirmed(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, trxId string) (bool, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return  false,utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmedRPC(trxId)
		return isConfirmed, err
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmed(trxId)
		return isConfirmed, err
	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmed(trxId)
		return isConfirmed, err
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isConfirmed, err := agent.IsTransactionConfirmedRPC(trxId)
		return isConfirmed, err
	}
	return false, nil
}

func IsAddressValid(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string) (bool, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return  false,utils.InvalidCoinSymbol(coinSymbol)
	}
	if address ==""{
		return true,nil
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isValid, err := agent.IsAddressValidRPC(address)
		return isValid, err
	}else if coinSymbol == "ETH"{
		isValied := ETHValidateAddress(address)
		return isValied, nil
	}else if coinObj.IsErc20{
		isValied := ETHValidateAddress(address)
		return isValied, nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		isValid, err := agent.IsAddressValidRPC(address)
		return isValid, err
	}
	return false, nil
}

func ConfirmTransaction(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string,
	trxId int, rawTrxId string) (error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB" || coinObj.IsOmni {
		trxMgr := model.GlobalDBMgr.TransactionMgr
		err := trxMgr.UpdateTransactionStateFeeCost(trxId, 2, nil)
		if err != nil {
			return err
		}
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		real_cost,err := agent.GetTransactionRealCost(coinSymbol,rawTrxId)
		if err != nil {
			return err
		}
		trxMgr := model.GlobalDBMgr.TransactionMgr
		err = trxMgr.UpdateTransactionStateFeeCost(trxId, 2, &real_cost)
		if err != nil {
			return err
		}
	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		real_cost,err := agent.GetTransactionRealCost(coinSymbol,rawTrxId)
		if err != nil {
			return err
		}
		trxMgr := model.GlobalDBMgr.TransactionMgr
		err = trxMgr.UpdateTransactionStateFeeCost(trxId, 2, &real_cost)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBalance(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, address string) (string,string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "","",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)
		balance, err := agent.GetBalanceByAddressRPC(address)
		if err != nil {
			return "","", err
		}

		return balance,"",nil
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		balance, err := agent.GetBalanceByAddress(address)
		if err != nil {
			return "","", err
		}

		return ConvertBigNumberToString(coinSymbol,balance),"", nil
	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip,rpcPort)
		agent.Init(url)
		erc20_balance,balance, err := agent.GetERC20BalanceByAddress(coinSymbol,address)
		if err != nil {
			return "","", err
		}

		return ConvertBigNumberToString(coinSymbol,erc20_balance),ConvertBigNumberToString("ETH",balance), nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		balance, err := agent.GetBalanceByAddressRPC(address)
		if err != nil {
			return "","", err
		}
		feeBalance, err := agent.GetFeeBalanceByAddressRPC(address)
		if err != nil {
			return "","", err
		}

		return balance,feeBalance,nil
	}
	return "","", nil
}

func Transfer(coinSymbol string, ip string, rpcPort int, rpcUser string, rpcPass string, keyIndex uint16, pubKeyStr string,
	from string, to string, amountStr string, feeStr string, gasPrice string, gasLimit string) (string, string, error) {
	coinObj,exist := config.GlobalSupportCoinMgr[coinSymbol]
	if !exist{
		return "","",utils.InvalidCoinSymbol(coinSymbol)
	}
	if coinSymbol == "BTC" || coinSymbol == "LTC" || coinSymbol == "BCH" || coinSymbol == "UB"{
		agent := AgentFactory(coinSymbol)
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// get utxos
		utxos, err := agent.GetUtxosByAddressRPC(from)
		if err != nil {
			return "", "", err
		}

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxInPutsOutPutsRPC(from, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}

		// sign transaction
		trxSigHex, err := agent.SignRawTransaction(rawTrx, pubKeyStr, keyIndex, utxos)
		if err != nil {
			return "", "", err
		}

		fmt.Println("trxSigHex:", trxSigHex)

		// broadcast transaction
		trxId, err := agent.BroadcastTransactionRPC(trxSigHex)
		if err != nil {
			return "", "", err
		}

		feeCostStr := feeCost
		//feeCostStr := fmt.Sprintf("%.08f", feeCost)
		return feeCostStr, trxId, nil
	}else if coinSymbol == "ETH"{
		var agent ETHAgent
		url := fmt.Sprintf("http://%s:%d",ip, rpcPort)
		agent.Init(url)
		rawTrx,err := agent.CreateTransaction(coinSymbol,from,to,gasLimit,gasPrice,amountStr,"",keyIndex)
		if err != nil {
			return "", "", err
		}
		trxId,err := agent.BroadcastTransaction(rawTrx)
		if err != nil {
			return "", "", err
		}
		return "", trxId, nil

	}else if coinObj.IsErc20{
		var agent ERC20Agent
		url := fmt.Sprintf("http://%s:%d",ip, rpcPort)
		agent.Init(url)
		rawTrx,err := agent.CreateERC20Transaction(coinSymbol,from,to,gasLimit,gasPrice,amountStr,"",keyIndex)
		if err != nil {
			return "", "", err
		}
		trxId,err := agent.BroadcastTransaction(rawTrx)
		if err != nil {
			return "", "", err
		}
		return "", trxId, nil
	} else if coinObj.IsOmni {
		agent := new(OMNIAgent)
		agent.CoinSymbol = coinObj.CoinSymbol
		url := fmt.Sprintf("http://%s:%s@%s:%d", rpcUser, rpcPass, ip, rpcPort)
		agent.Init(url)

		// get utxos
		utxos, err := agent.GetUtxosByAddressRPC(from)
		if err != nil {
			return "", "", err
		}

		// create raw transaction
		feeCost, inPuts, outPuts, err := agent.BuildTrxInPutsOutPutsRPC(from, to, amountStr, feeStr)
		if err != nil {
			return "", "", err
		}
		rawTrx, err := agent.CreateRawTransactionRPC(inPuts, outPuts)
		if err != nil {
			return "", "", err
		}
		rawTrx, err = agent.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, uint32(coinObj.OmniPropertyId), amountStr)
		if err != nil {
			return "", "", err
		}

		// sign transaction
		trxSigHex, err := agent.SignRawTransaction(rawTrx, pubKeyStr, keyIndex, utxos)
		if err != nil {
			return "", "", err
		}

		fmt.Println("trxSigHex:", trxSigHex)

		// broadcast transaction
		trxId, err := agent.BroadcastTransactionRPC(trxSigHex)
		if err != nil {
			return "", "", err
		}
		feeCostStr := feeCost
		//feeCostStr := fmt.Sprintf("%.08f", feeCost)
		return feeCostStr, trxId, nil
	}
	return "", "", nil
}

// sigType: '1' r+s,  '2' der
func CoinSignTrx(sigType byte, signData []byte, keyIndex uint16) ([]byte, error) {
	conn := new(buffer_tcp.BufferTcpConn)
	err := conn.TCPConnect(config.GlobalConfig.CryptoDeviceConfig.DeviceIp,
		config.GlobalConfig.CryptoDeviceConfig.DevicePort, float64(config.GlobalConfig.CryptoDeviceConfig.TimeOut))
	if err != nil {
		return nil, err
	}
	defer conn.TCPDisConnect()

	var l7req cryptocard.L7Request
	l7req.Set(sigType, keyIndex, nil, signData)
	err = l7req.Pack(conn)
	if err != nil {
		return nil, err
	}

	var l8resp cryptocard.L8Response
	err = l8resp.UnPack(conn)
	if err != nil {
		return nil, err
	}

	return l8resp.DataSigned, nil
}

// sigType: '1' r+s,  '2' der
func CoinVerifyTrx(sigType byte, keyIndex uint16, signData []byte, signedData []byte) (bool, error){
	conn := new(buffer_tcp.BufferTcpConn)
	err := conn.TCPConnect(config.GlobalConfig.CryptoDeviceConfig.DeviceIp,
		config.GlobalConfig.CryptoDeviceConfig.DevicePort, float64(config.GlobalConfig.CryptoDeviceConfig.TimeOut))
	if err != nil {
		return false, err
	}
	defer conn.TCPDisConnect()

	var l4req cryptocard.L4Request
	l4req.Set(sigType, keyIndex, nil, nil, signData, signedData)
	err = l4req.Pack(conn)
	if err != nil {
		return false, err
	}

	var l5resp cryptocard.L5Response
	err = l5resp.UnPack(conn)
	if err != nil {
		return false, err
	}

	if l5resp.ErrCode[0] == '0' && l5resp.ErrCode[1] == '0' {
		return true, nil
	} else {
		return false, nil
	}

	return true, nil
}

func SerializeDerEncoding(rBytes []byte, sBytes []byte) ([]byte, error) {
	if len(rBytes) != 32 {
		return nil, errors.New("invalid rBytes len")
	}
	if len(sBytes) != 32 {
		return nil, errors.New("invalid sBytes len")
	}

	var r []byte
	r = append(r, 0)
	r = append(r, rBytes...)
	var s []byte
	s = append(s, 0)
	s = append(s, sBytes...)
	for {
		if len(r) > 1 && r[0] == 0 && r[1] < 0x80 {
			r = r[1:]
		} else {
			break
		}
	}
	for {
		if len(s) > 1 && s[0] == 0 && s[1] < 0x80 {
			s = s[1:]
		} else {
			break
		}
	}

	size := 6 + len(r) + len(s)
	signedData := make([]byte, size, size)
	signedData[0] = 0x30
	signedData[1] = 4 + byte(len(r)) + byte(len(s))
	signedData[2] = 0x2
	signedData[3] = byte(len(r))
	copy(signedData[4:4+len(r)], r)
	signedData[4+len(r)] = 0x2
	signedData[5+len(r)] = byte(len(s))
	copy(signedData[6+len(r):6+len(r)+len(s)], s)

	return signedData, nil
}

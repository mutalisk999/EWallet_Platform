package coin

import (
	"fmt"
	"testing"
	"config"
)

func TestBCHAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	ba, err := ag.GetBalanceByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestBCHAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	c, err := ag.IsTransactionConfirmedRPC("4a6f14efecc0f10a8afff7dd1a9d1c2b69e2c1588ab512fd85b023f0c19fcc1f")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBCHAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10003")
	c, err := ag.IsAddressValidRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBCHAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestBCHCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "1ccd06a58246e58f58e339940ad1a994a528cae4bf4c43f2eeddc1cc245779c6cb1674ba3bfaf36f5a2502981d9ce158eaca4d1057a61996f2c755a02bdc6e03"
	config.IsTestEnvironment = true
	addrStr, _ := BCHCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestBCHAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestBCHAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestBCHUnPackRawTransaction(t *testing.T) {
	rawTrx := "02000000011fcc9fc1f023b085fd12b58a58c1e2692b1c9d1addf7ff8a0af1c0ecef146f4a0100000000ffffffff0240420f00000000001976a9146180cf77dca5a0a5f761e9cef0b1a5c917384e8088ac18a4eb02000000001976a9149e6827242840668d9bb217498c3a29207f9ed44888ac00000000"
	trx, _ := BCHUnPackRawTransaction(rawTrx)
	fmt.Println("trx version:", trx.Version)
	fmt.Println("trx locktime", trx.LockTime)
	fmt.Println("trx vin size:", len(trx.Vin))
	for i := 0; i < len(trx.Vin); i++ {
		fmt.Println("vin prevout:", trx.Vin[i].PrevOut.Hash.GetHex(), trx.Vin[i].PrevOut.N)
		fmt.Println("vin scriptsig:", trx.Vin[i].ScriptSig)
		fmt.Println("vin sequence:", trx.Vin[i].Sequence)
		fmt.Println("vin scriptwitness:", trx.Vin[i].ScriptWitness)
	}
}

func TestBCHSignRawTransaction(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "1ccd06a58246e58f58e339940ad1a994a528cae4bf4c43f2eeddc1cc245779c6cb1674ba3bfaf36f5a2502981d9ce158eaca4d1057a61996f2c755a02bdc6e03"
	keyIndex := uint16(4)
	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1
	trxSigHex, err := ag.SignRawTransaction(rawTrx, pubKeyStr, keyIndex, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSigHex:", trxSigHex)
}

func TestBCHAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("BCH")
	ag.Init("http://test:test@192.168.1.124:10003")
	utxos, err := ag.GetUtxosByAddressRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("bchtest:qz0xsfey9pqxdrvmkgt5nrp69ys8l8k5fqgmz25607", "bchtest:qpscpnmhmjj6pf0hv85uau935hy3wwzwsqswd643qk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "1ccd06a58246e58f58e339940ad1a994a528cae4bf4c43f2eeddc1cc245779c6cb1674ba3bfaf36f5a2502981d9ce158eaca4d1057a61996f2c755a02bdc6e03"
	keyIndex := uint16(4)

	config.GlobalConfig.CryptoDeviceConfig.DeviceIp = "192.168.1.188"
	config.GlobalConfig.CryptoDeviceConfig.DevicePort = 1818
	config.GlobalConfig.CryptoDeviceConfig.TimeOut = 1
	trxSigHex, err := ag.SignRawTransaction(rawTrx, pubKeyStr, keyIndex, utxos)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("trxSigHex:", trxSigHex)
	trxId, err := ag.BroadcastTransactionRPC(trxSigHex)
	if err == nil {
		fmt.Println("trxId:", trxId)
	} else {
		fmt.Println("err:", err)
	}
}


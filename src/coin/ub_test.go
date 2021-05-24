package coin

import (
	"fmt"
	"testing"
	"config"
)

func TestUBAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	ba, err := ag.GetBalanceByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestUBAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	c, err := ag.IsTransactionConfirmedRPC("02faec58bfbcb47a34b4670b25fe2117d69e8aa234eeee2d62b262a75fd63d65")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestUBAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	c, err := ag.IsAddressValidRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestUBAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestUBCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "9298d576117276c4eb10adcb1af26c2a7779b15aa2cae9941512b4732e7ed022ea0bc212a0809f86c0da620b2f5ea37befcff7e253ca5f714a2c68e632b3522e"
	config.IsTestEnvironment = true
	addrStr, _ := UBCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestUBAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestUBAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestUBUnPackRawTransaction(t *testing.T) {
	rawTrx := "02000000017d704bb25c19e317c2db09b95d364a57ecd4a81b3584f647dd7e4dbf0599e9060100000000ffffffff0240420f00000000001976a914130c91e9def87a4445440ab430f7a98e76424d9188ac98472677000000001976a9141b46aa5c903f3dc8eb9592876af0a061db4b3bed88ac00000000"
	trx, _ := UBUnPackRawTransaction(rawTrx)
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

func TestUBSignRawTransaction(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "9298d576117276c4eb10adcb1af26c2a7779b15aa2cae9941512b4732e7ed022ea0bc212a0809f86c0da620b2f5ea37befcff7e253ca5f714a2c68e632b3522e"
	keyIndex := uint16(6)

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

func TestUBAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("UB")
	ag.Init("http://test:test@192.168.1.124:10004")
	utxos, err := ag.GetUtxosByAddressRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mi1BCA3Skdv4jAcyPejebDFrzCC3uELYUn", "mhFgC5MnmkYQKMrpDHiTfGdrPqgHkYpy4U", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "9298d576117276c4eb10adcb1af26c2a7779b15aa2cae9941512b4732e7ed022ea0bc212a0809f86c0da620b2f5ea37befcff7e253ca5f714a2c68e632b3522e"
	keyIndex := uint16(6)
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


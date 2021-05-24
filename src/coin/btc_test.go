package coin

import (
	"fmt"
	"testing"
	"config"
)

func TestBTCAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	ba, err := ag.GetBalanceByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestBTCAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	c, err := ag.IsTransactionConfirmedRPC("ed410a120a23c9c8e078f2ed8f43aa99f05bd95090e61befee0de00644e578ee")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBTCAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	c, err := ag.IsAddressValidRPC("mhuXoAkUNLPcboTFu9PDtGapc3wnZAGeyw")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestBTCAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestBTCCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f"
	config.IsTestEnvironment = true
	addrStr, _ := BTCCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestBTCAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestBTCAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestBTCUnPackRawTransaction(t *testing.T) {
	rawTrx := "02000000016059331a0741e2be764f25f7c85fa78855aa63cdc780ca10bfcac168ad7fd7ad0000000000ffffffff0240420f000000000017a914f1b3b098ae94b096b60b6bfc04a51094ede15a2887184a8900000000001976a91430e83ac38de345ce80f069055cb9f8e15b28e54d88ac00000000"
	trx, _ := BTCUnPackRawTransaction(rawTrx)
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

func TestBTCSignRawTransaction(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f"
	keyIndex := uint16(1)

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

func TestBTCAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("BTC")
	ag.Init("http://test:test@192.168.1.124:10001")
	utxos, err := ag.GetUtxosByAddressRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjyYvrTuYRGYoHqowFHqwriuiaFfcRWFp7", "2NFHE6anahifG7o7dhkrseNHaFBJ2x53YDk", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f"
	keyIndex := uint16(1)
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


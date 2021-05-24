package coin

import (
	"fmt"
	"testing"
	"config"
)

func TestLTCAgent_GetBalanceByAddress(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	ba, err := ag.GetBalanceByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestLTCAgent_IsTransactionConfirmed(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	c, err := ag.IsTransactionConfirmedRPC("9d5675ffe98872f1fd02b548032de6e40a48259430ae8c00f185df883ff18d81")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestLTCAgent_IsAddressValidRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	c, err := ag.IsAddressValidRPC("n1uaazxSzWochahoAnKbGPAxh34MhxNe1J")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestLTCAgent_GetUtxosByAddress(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestLTCCalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "4d12208801f9cfc25ff0cb62afb6affef4e636cc2306f46987f16366b78807010aa542dac5fb971f25928b1a2d53267c00358ddf590c203c865bfb2aa88eb17a"
	config.IsTestEnvironment = true
	addrStr, _ := LTCCalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestLTCAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestLTCAgent_CreateRawTransaction(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestLTCUnPackRawTransaction(t *testing.T) {
	rawTrx := "0200000001dd31ca8a5c4417867072f7c795e482b95677e942d1b81db00e50fe60ac19a8de0100000000ffffffff0240420f00000000001976a91404b80434d66065d51639e05ac0af1b885058b3de88ac30941201000000001976a914ab9ffaedde51a4b4720ee51208b568d1834b68ef88ac00000000"
	trx, _ := LTCUnPackRawTransaction(rawTrx)
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

func TestLTCSignRawTransaction(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "4d12208801f9cfc25ff0cb62afb6affef4e636cc2306f46987f16366b78807010aa542dac5fb971f25928b1a2d53267c00358ddf590c203c865bfb2aa88eb17a"
	keyIndex := uint16(3)
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

func TestLTCAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := AgentFactory("LTC")
	ag.Init("http://test:test@192.168.1.124:10002")
	utxos, err := ag.GetUtxosByAddressRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mwARVUYJUYJRDhxZb7YoAAg4D46Zhv8Ngh", "mfwuSWNFQEHeCYjVcHrmm5a5H3GJSDsB34", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	pubKeyStr := "4d12208801f9cfc25ff0cb62afb6affef4e636cc2306f46987f16366b78807010aa542dac5fb971f25928b1a2d53267c00358ddf590c203c865bfb2aa88eb17a"
	keyIndex := uint16(3)

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


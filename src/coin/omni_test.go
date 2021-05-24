package coin

import (
	"fmt"
	"testing"
	"config"
)

func TestOMNIAgent_GetFeeBalanceByAddress(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	ba, err := ag.GetFeeBalanceByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestOMNIAgent_GetBalanceByAddress(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	ba, err := ag.GetBalanceByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(ba)
}

func TestOMNIAgent_IsFeeTransactionConfirmed(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	c, err := ag.IsFeeTransactionConfirmedRPC("08d75f23070a12bb1059e3498dab47a000ce303238f99584da412b6a4c6b9558")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestOMNIAgent_IsTransactionConfirmed(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	c, err := ag.IsTransactionConfirmedRPC("362ec25c3f877f016918babca120cf80bfbc05051091e7079bb16b8ddf5807cf")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestOMNIAgent_IsAddressValidRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	c, err := ag.IsAddressValidRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(c)
}

func TestOMNIAgent_GetUtxosByAddress(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println(utxos)
}

func TestOMNICalcAddressByPubKey(t *testing.T) {
	pubKeyStr := "b389cf01cacc2aae68942d2c218b40dbb31cce668d61276dfdf6f823c91daeb8046b6b6f5e0716a2882c38a140d417c4f049c4c11a2799b25b6ded66302ad738"
	config.IsTestEnvironment = true
	addrStr, _ := OMNICalcAddressByPubKey(pubKeyStr)
	fmt.Println(addrStr)
}

func TestOMNIAgent_BuildTrxInPutsOutPuts(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	feeCost, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	fmt.Println(feeCost)
	fmt.Println(inPuts)
	fmt.Println(outPuts)
}

func TestOMNIAgent_CreateRawTransaction(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	fmt.Println(rawTrx)
}

func TestOMNIAgent_CreateRawTransactionOpReturnRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	rawTrx := "0100000001f1238a114f4a2ca83a5dbbc04832dcf0e8d8e4c02eeaecab22f03f24bb9be5e50100000000ffffffff02a6889800000000001976a9142996789a3cdd905d8e44fb6033a08c667410ef6e88ac22020000000000001976a9149f2a66d7b349ba9a87e6bf9cf7da1df31697d63988ac00000000"
	rawTrx, _  = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	fmt.Println(rawTrx)
}

func TestOMNIUnPackRawTransaction(t *testing.T) {
	rawTrx := "0100000001f1238a114f4a2ca83a5dbbc04832dcf0e8d8e4c02eeaecab22f03f24bb9be5e50100000000ffffffff03a6889800000000001976a9142996789a3cdd905d8e44fb6033a08c667410ef6e88ac22020000000000001976a9149f2a66d7b349ba9a87e6bf9cf7da1df31697d63988ac0000000000000000166a146f6d6e69000000000000000200000000000f424000000000"
	trx, _ := OMNIUnPackRawTransaction(rawTrx)
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

func TestOMNISignRawTransaction(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	rawTrx, _  = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	pubKeyStr := "b389cf01cacc2aae68942d2c218b40dbb31cce668d61276dfdf6f823c91daeb8046b6b6f5e0716a2882c38a140d417c4f049c4c11a2799b25b6ded66302ad738"
	keyIndex := uint16(5)

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

func TestOMNIAgent_BroadcastTransactionRPC(t *testing.T) {
	ag := new(OMNIAgent)
	ag.CoinSymbol = "TOMNI"
	ag.Init("http://test:test@192.168.1.124:10009")
	utxos, err := ag.GetUtxosByAddressRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	_, inPuts, outPuts, _ := ag.BuildTrxInPutsOutPutsRPC("mjJrFqMw2u2jfssXJf8PMoS6EMibC1zng6", "mv2YXgKpgVqaaus6zdzJGtrWEQ4iPBXyvV", "0.01", "0.0001")
	rawTrx, _ := ag.CreateRawTransactionRPC(inPuts, outPuts)
	rawTrx, _  = ag.CreateRawTransactionOpReturnRPC(rawTrx, 0, 0, 2, "0.01")
	pubKeyStr := "b389cf01cacc2aae68942d2c218b40dbb31cce668d61276dfdf6f823c91daeb8046b6b6f5e0716a2882c38a140d417c4f049c4c11a2799b25b6ded66302ad738"
	keyIndex := uint16(5)
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


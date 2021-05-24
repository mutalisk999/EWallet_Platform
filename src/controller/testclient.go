package controller

import (
	"fmt"
	"github.com/ybbus/jsonrpc"
	"strconv"
	"utils"
)

type account struct {
	pubkey string
	prikey string
}
type testClient struct {
	url string
}

func (agent *testClient) doHttpJsonRpcCallType1(appdenurl string, method string, args interface{}) (*jsonrpc.RPCResponse, error) {
	rpcClient := jsonrpc.NewClient(agent.url + appdenurl)
	rpcResponse, err := rpcClient.Call(method, args)
	if err != nil {
		return nil, err
	}
	return rpcResponse, nil
}

//return authid,err
func (c *testClient) get_authid() (int64, error) {
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = ""
	pa0["idtype"] = 1
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/identity", "get_auto_inc_id", pa)
	if err != nil {
		return 0, err
	}
	return res.GetInt()
}

//return mgmtid,err
func (c *testClient) get_mgmtid(sid string, idtp int) (int64, error) {
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["sessionid"] = sid
	pa0["idtype"] = idtp
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/identity", "get_auto_inc_id", pa)
	if err != nil {
		return 0, err
	}
	return res.GetInt()
}

//return session_id,err
func (c *testClient) login(auth_id int64, acc account) (string, error) {
	pa := make([]map[string]interface{}, 0)
	var pa0 map[string]interface{}
	pa0 = make(map[string]interface{})
	pa0["loginid"] = auth_id
	pa0["pubkey"] = acc.pubkey
	sigdata := "user_login," + strconv.FormatInt(auth_id, 10)
	sig_res, err := utils.RsaSignWithSha1Hex(sigdata, acc.prikey)
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(sig_res)
	pa0["signature"] = sig_res
	pa = append(pa, pa0)
	res, err := c.doHttpJsonRpcCallType1("/apis/user", "user_login", pa)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		println(res.Error.Message)
		return "", nil
	}
	return res.Result.(map[string]interface{})["sessionid"].(string), nil
}

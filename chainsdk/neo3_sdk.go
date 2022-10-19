/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package chainsdk

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/nep17"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/joeqian10/neo3-gogogo/rpc/models"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type Neo3Sdk struct {
	client *rpc.RpcClient
	url    string
}

type Neo3RpcReq struct {
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      uint     `json:"id"`
}

type Nep11Property2 struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Name      string `json:"name"`
		Image     string `json:"image"`
		Series    string `json:"series"`
		Supply    string `json:"supply"`
		Thumbnail string `json:"thumbnail"`
	} `json:"result"`
}

type Nep11Property struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Series    string `json:"series"`
	Supply    string `json:"supply"`
	Thumbnail string `json:"thumbnail"`
}

func NewNeo3Sdk(url string) *Neo3Sdk {
	return &Neo3Sdk{
		client: rpc.NewClient(url),
		url:    url,
	}
}

func (sdk *Neo3Sdk) GetClient() *rpc.RpcClient {
	return sdk.client
}

func (sdk *Neo3Sdk) GetUrl() string {
	return sdk.url
}

func (sdk *Neo3Sdk) GetBlockCount() (uint64, error) {
	res := sdk.client.GetBlockCount()
	if res.ErrorResponse.Error.Message != "" {
		return 0, fmt.Errorf("%s", res.ErrorResponse.Error.Message)
	}
	return uint64(res.Result), nil
}

func (sdk *Neo3Sdk) GetBlockByIndex(index uint64) (*models.RpcBlock, error) {
	res := sdk.client.GetBlock(strconv.Itoa(int(index)))
	if res.ErrorResponse.Error.Message != "" {
		return nil, fmt.Errorf("%s", res.ErrorResponse.Error.Message)
	}
	return &res.Result, nil
}

func (sdk *Neo3Sdk) GetApplicationLog(txId string) (*models.RpcApplicationLog, error) {
	res := sdk.client.GetApplicationLog(txId)
	if res.ErrorResponse.Error.Message != "" {
		return nil, fmt.Errorf("%s", res.ErrorResponse.Error.Message)
	}
	return &res.Result, nil
}

func (sdk *Neo3Sdk) GetTransactionHeight(hash string) (uint64, error) {
	res := sdk.client.GetTransactionHeight(hash)
	if res.ErrorResponse.Error.Message != "" {
		return 0, fmt.Errorf("%s", res.ErrorResponse.Error.Message)
	}
	return uint64(res.Result), nil
}

func (sdk *Neo3Sdk) SendRawTransaction(txHex string) (bool, error) {
	res := sdk.client.SendRawTransaction(txHex)
	if res.HasError() {
		return false, fmt.Errorf("%s", res.ErrorResponse.Error.Message)
	}
	return true, nil
}

func (sdk *Neo3Sdk) Nep17Info(hash string) (string, string, int64, error) {
	scriptHash, err := helper.UInt160FromString(hash)
	if err != nil {
		return "", "", 0, err
	}
	nep17 := nep17.NewNep17Helper(scriptHash, sdk.client)
	decimal, err := nep17.Decimals()
	if err != nil {
		return "", "", 0, err
	}
	symbol, err := nep17.Symbol()
	if err != nil {
		return "", "", 0, err
	}
	return hash, symbol, int64(decimal), nil
}

func Neo3AddrToHash160(addr string) (string, error) {
	scriptHash, err := crypto.AddressToScriptHash(addr, helper.DefaultAddressVersion)
	return scriptHash.String(), err
}

func Hash160ToNeo3Addr(encodedHash string) (string, error) {
	decodedByte, err := crypto.Base64Decode(encodedHash)
	if err != nil {
		return "", err
	}
	hash160 := helper.UInt160FromBytes(decodedByte)
	return crypto.ScriptHashToAddress(hash160, helper.DefaultAddressVersion), nil
}

func (sdk *Neo3Sdk) Nep11Property() (*Nep11Property, error) {
	asset := "0x4fb2f93b37ff47c0c5d14cfc52087e3ca338bc56"
	tokenId := "4d65746150616e616365612023302d3031"
	params := make([]string, 2)
	params[0] = asset
	params[1] = tokenId
	body, err := jsonRequest(sdk.url, "getnep11properties", params)
	if err != nil {
		return nil, err
	}
	var resp *Nep11Property
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sdk *Neo3Sdk) ContractCall1() {
	hash := "0x4fb2f93b37ff47c0c5d14cfc52087e3ca338bc56"
	method := "tokensOf"
	var params []models.RpcContractParameter
	addr := "Nd6UWMyUDp1nZK7osMXJW9s21NqDNmnBfz"
	addrHash, _ := Neo3AddrToHash160(addr)
	params = append(params, models.RpcContractParameter{
		Type:  "Hash160",
		Value: addrHash,
	})
	res := sdk.client.InvokeFunction(hash, method, params, nil)
	aaa, _ := json.MarshalIndent(res, "", "	")
	fmt.Println(string(aaa))
}

func (sdk *Neo3Sdk) ContractCall2() {

	//hash1 := "vHbGEieTw0YQMJtHu+ELwwRiRVg="
	//fmt.Println(Hash160ToNeo3Addr(hash1))

	fmt.Println()

	//fmt.Println(Hash160ToNeo3Addr(stack.Value))

	//p, err := stack.ToParameter()
	//
	//fmt.Println(string(p.Value.([]byte)))
	//val, ok := stack.Value.(map[models.InvokeStack]models.InvokeStack)
	//if ok {
	//	for k, v := range val {
	//		res, _ := k.ToParameter()
	//		ss := res.Value.([]byte)
	//		fmt.Println(string(ss))
	//		res2, _ := v.ToParameter()
	//		ss2 := res2.Value.([]byte)
	//		fmt.Println(string(ss2))
	//	}
	//}
}

func (sdk *Neo3Sdk) Nep11OwnerOf(assetHash, tokenId string) (string, error) {
	method := "ownerOf"
	var params []models.RpcContractParameter
	tokenIdBase64 := crypto.Base64Encode(helper.HexToBytes(tokenId))
	params = append(params, models.RpcContractParameter{
		Type:  "ByteArray",
		Value: tokenIdBase64,
	})
	response := sdk.client.InvokeFunction(assetHash, method, params, nil)
	stack, err := rpc.PopInvokeStack(response)
	if err != nil {
		return "", err
	}
	return Hash160ToNeo3Addr(stack.Value.(string))
}
func (sdk *Neo3Sdk) Nep11BalanceOf(assetHash, owner string) (*big.Int, error) {
	method := "balanceOf"
	ownerHash160, _ := Neo3AddrToHash160(owner)
	var params []models.RpcContractParameter
	params = append(params, models.RpcContractParameter{
		Type:  "Hash160",
		Value: ownerHash160,
	})
	response := sdk.client.InvokeFunction(assetHash, method, params, nil)
	stack, err := rpc.PopInvokeStack(response)
	if err != nil {
		return nil, err
	}
	val, _ := stack.ToParameter()
	return val.Value.(*big.Int), nil
}

func (sdk *Neo3Sdk) Nep11TokensOf(assetHash, owner string) ([]string, error) {
	method := "tokensOf"
	ownerHash160, _ := Neo3AddrToHash160(owner)
	var params []models.RpcContractParameter
	params = append(params, models.RpcContractParameter{
		Type:  "Hash160",
		Value: ownerHash160,
	})
	response := sdk.client.InvokeFunction(assetHash, method, params, nil)
	stack, err := rpc.PopInvokeStack(response)
	if err != nil {
		return nil, err
	}
	val, _ := stack.ToParameter()
	return val.Value.([]string), nil
}

func (sdk *Neo3Sdk) Nep11Properties(assetHash, tokenId string) (*Nep11Property, error) {
	method := "properties"
	var params []models.RpcContractParameter
	tokenIdBase64 := crypto.Base64Encode(helper.HexToBytes(tokenId))
	params = append(params, models.RpcContractParameter{
		Type:  "ByteArray",
		Value: tokenIdBase64,
	})
	response := sdk.client.InvokeFunction(assetHash, method, params, nil)
	stack, err := rpc.PopInvokeStack(response)
	if err != nil {
		return nil, err
	}
	property := &Nep11Property{}
	val, ok := stack.Value.(map[models.InvokeStack]models.InvokeStack)
	if ok {
		propertyMap := make(map[string]string)
		var propertyKey, propertyVal string
		for k, v := range val {
			res, _ := k.ToParameter()
			propertyKey = string(res.Value.([]byte))
			res2, _ := v.ToParameter()
			propertyVal = string(res2.Value.([]byte))
			propertyMap[propertyKey] = propertyVal
		}
		arr, _ := json.Marshal(propertyMap)
		_ = json.Unmarshal(arr, &property)
	}
	return property, nil
}

func (sdk *Neo3Sdk) Nep11TokenUrl(assetHash, tokenId string) (string, error) {
	property, err := sdk.Nep11Properties(assetHash, tokenId)
	return property.Image, err
}

func (sdk *Neo3Sdk) Nep17Balance(hash string, addr string) (*big.Int, error) {
	scriptHash, err := helper.UInt160FromString(hash)
	if err != nil {
		return new(big.Int).SetUint64(0), err
	}
	nep17 := nep17.NewNep17Helper(scriptHash, sdk.client)
	addrHash, err := helper.UInt160FromString(addr)
	if err != nil {
		logs.Info("Nep17Balance err: %s", err)
		return new(big.Int).SetUint64(0), err
	}
	logs.Info("Nep17Balance addrHash: %+v", addrHash)
	return nep17.BalanceOf(addrHash)
}

func (sdk *Neo3Sdk) Nep17TotalSupply(hash string) (*big.Int, error) {
	scriptHash, err := helper.UInt160FromString(hash)
	if err != nil {
		return new(big.Int).SetUint64(0), err
	}
	logs.Info("hash: %s", hash)
	nep17 := nep17.NewNep17Helper(scriptHash, sdk.client)
	if err != nil {
		return new(big.Int).SetUint64(0), err
	}
	return nep17.TotalSupply()
}

func jsonRequest(url, method string, params []string) (result []byte, err error) {
	req := &Neo3RpcReq{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	data, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

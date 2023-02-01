package test

import (
	"encoding/json"
	"fmt"
	"os"
	"poly-bridge/basedef"
	"poly-bridge/conf"
	"poly-bridge/crosschainlisten"
	"testing"
)

func TestEthereumChainListen_HandleNewBatchBlock(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./config.json")
	if config == nil {
		panic("read config failed!")
	}
	ethListenConfig := config.GetChainListenConfig(basedef.ZION_CROSSCHAIN_ID)
	if ethListenConfig == nil {
		panic("config is not valid")
	}
	chainHandle := crosschainlisten.NewChainHandle(ethListenConfig)
	_, _, polyTransactions, _, _, _, err := chainHandle.HandleNewBatchBlock(1, 500)
	if err != nil {
		fmt.Println("err", err)
	}
	a, _ := json.MarshalIndent(polyTransactions, "", "	")
	fmt.Println(string(a))
}

func TestEthereumChainListen_HandleNewBatchBlock2(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./config.json")
	if config == nil {
		panic("read config failed!")
	}
	ethListenConfig := config.GetChainListenConfig(basedef.GOERLI_CROSSCHAIN_ID)
	if ethListenConfig == nil {
		panic("config is not valid")
	}
	chainHandle := crosschainlisten.NewChainHandle(ethListenConfig)
	_, _, _, dstTransactions, _, _, err := chainHandle.HandleNewBatchBlock(8408230, 8409230)
	if err != nil {
		fmt.Println("err", err)
	}
	a, _ := json.MarshalIndent(dstTransactions, "", "	")
	fmt.Println(string(a))
}

func TestEthereumChainListen_HandleNewBlock3(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./config.json")
	if config == nil {
		panic("read config failed!")
	}
	ethListenConfig := config.GetChainListenConfig(basedef.ONT_CROSSCHAIN_ID)
	if ethListenConfig == nil {
		panic("config is not valid")
	}
	chainHandle := crosschainlisten.NewChainHandle(ethListenConfig)
	a, b, c, d, _, _, _, _, err := chainHandle.HandleNewBlock(17214676)
	if err != nil {
		fmt.Println("err", err)
	}
	e, _ := json.MarshalIndent(a, "", "	")
	fmt.Println("a", string(e))
	e, _ = json.MarshalIndent(b, "", "	")
	fmt.Println("b", string(e))
	e, _ = json.MarshalIndent(c, "", "	")
	fmt.Println("c", string(e))
	e, _ = json.MarshalIndent(d, "", "	")
	fmt.Println("d", string(e))
}

package chainsdk

import (
	"fmt"
	"testing"
)

func TestNeo3Sdk_ContractCall1(t *testing.T) {
	sdk := NewNeo3Sdk("http://seed1t5.neo.org:20332")
	sdk.ContractCall1()
}

func TestNeo3Sdk_ContractCall2(t *testing.T) {
	sdk := NewNeo3Sdk("http://seed1t5.neo.org:20332")
	sdk.ContractCall2()
}

func TestNeo3Sdk_Nep11Property(t *testing.T) {
	sdk := NewNeo3Sdk("http://seed1t5.neo.org:20332")
	fmt.Println(sdk.Nep11OwnerOf("0x4fb2f93b37ff47c0c5d14cfc52087e3ca338bc56", "4d65746150616e616365612023302d3031"))
}

func TestNeo3Sdk_Nep11Properties(t *testing.T) {
	sdk := NewNeo3Sdk("http://seed1t5.neo.org:20332")
	fmt.Println(sdk.Nep11Properties("0x4fb2f93b37ff47c0c5d14cfc52087e3ca338bc56", "4d65746150616e616365612023302d3031"))
}

func TestNeo3Sdk_Nep11BalanceOf(t *testing.T) {
	sdk := NewNeo3Sdk("http://seed1t5.neo.org:20332")
	fmt.Println(sdk.Nep11BalanceOf("0x4fb2f93b37ff47c0c5d14cfc52087e3ca338bc56", "Nd6UWMyUDp1nZK7osMXJW9s21NqDNmnBfz"))

}

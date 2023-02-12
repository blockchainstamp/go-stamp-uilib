package main

import (
	"flag"
	"fmt"
	contract "github.com/blockchainstamp/go-stamp-wallet/smart_contract"
	"testing"
)

var (
	addr = ""
)

func init() {
	flag.StringVar(&addr, "addr", "0xF9Cbfd74808f812a3B8A2337BFC426B2A10Bd74a", "--addr")
}
func TestStampBalanceOfWallet(t *testing.T) {
}

func TestStampConfig(t *testing.T) {
	fmt.Println(contract.StampConfigFromBlockChain(addr))
}

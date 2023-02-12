package main

import "C"
import (
	"encoding/json"
	"fmt"
	contract "github.com/blockchainstamp/go-stamp-wallet/smart_contract"
	"github.com/ethereum/go-ethereum/common"
)

//export GetBalance
func GetBalance(wAddr, sAddr string) *C.char {
	r, e := contract.StampBalanceOfWallet(common.HexToAddress(sAddr), common.HexToAddress(wAddr))
	if e != nil {
		_appInst.SetError(e.Error())
		return nil
	}
	bts, _ := json.Marshal(r)
	return C.CString(string(bts))
}

//export StampConfig
func StampConfig(sAddr string) *C.char {

	fmt.Println("======>>>", sAddr)
	c, e := contract.StampConfigFromBlockChain(sAddr)
	if e != nil {
		_appInst.SetError(e.Error())
		return nil
	}
	fmt.Println("======>>>config:\n", c.String())
	return C.CString(c.String())
}

//export IsValidStampAddr
func IsValidStampAddr(sAddr string) bool {
	return common.IsHexAddress(sAddr)
}

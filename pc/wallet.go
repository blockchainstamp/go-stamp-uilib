package main

import "C"
import (
	"encoding/json"
	"fmt"
	bStamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
)

//export OpenWallet
func OpenWallet(addr, auth string) bool {
	_, err := bStamp.Inst().ActiveWallet(comm.WalletAddr(addr), auth)
	if err != nil {
		fmt.Println("======>>> active wallet err:", err)
		_appInst.SetError(err.Error())
		return false
	}

	return true
}

//export CreateWallet
func CreateWallet(auth, name string) *C.char {
	w, e := bStamp.Inst().CreateWallet(auth, name)
	if e != nil {
		_appInst.SetError(e.Error())
		return nil
	}
	return C.CString(w.String())
}

//export ImportWallet
func ImportWallet(wStr, auth string) *C.char {
	w, e := bStamp.Inst().ImportWallet(wStr, auth)
	if e != nil {
		_appInst.SetError(e.Error())
		return nil
	}
	return C.CString(w.String())
}

type WalletInfos struct {
	Addr    string
	Name    string
	JsonStr string
}

//export AllWallets
func AllWallets() *C.char {
	data := bStamp.Inst().ListAllWalletAddr()
	if len(data) == 0 {
		return C.CString("")
	}
	value := make(map[string]struct{}, 0)
	err := json.Unmarshal([]byte(data), &value)
	if err != nil {
		_appInst.SetError(err.Error())
		return nil
	}
	result := make([]*WalletInfos, 0)
	for addr, _ := range value {
		w, e := bStamp.Inst().GetWallet(comm.WalletAddr(addr))
		if e != nil {
			_appInst.SetError(err.Error())
			return nil
		}
		wi := &WalletInfos{
			Addr:    addr,
			Name:    w.Name(),
			JsonStr: w.Verbose(),
		}
		result = append(result, wi)
	}
	bts, err := json.Marshal(result)
	if err != nil {
		_appInst.SetError(err.Error())
		return nil
	}
	return C.CString(string(bts))
}

//export RemoveWallet
func RemoveWallet(wStr string) bool {
	e := bStamp.Inst().RemoveWallet(comm.WalletAddr(wStr))
	if e != nil {
		fmt.Println("======>>>remove failed", wStr, e.Error())
		_appInst.SetError(e.Error())
		return false
	}
	fmt.Println("======>>>remove success", wStr)
	return true
}

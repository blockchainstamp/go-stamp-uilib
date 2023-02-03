package main

/*
#include "callback.h"
*/
import "C"
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bStamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/sirupsen/logrus"
	"path"
)

func main() {
}

type App struct {
	callback C.UserInterfaceAPI
	setErr   C.SetLastErr
	localTls *tls.Config
}

var _appInst = &App{}

//export InitLib
func InitLib(baseDir, logLevel string, cb C.UserInterfaceAPI, errSet C.SetLastErr) bool {
	_appInst.callback = cb
	_appInst.setErr = errSet

	if err := bStamp.InitSDK(baseDir); err != nil {
		_appInst.SetError(err.Error())
		return false
	}
	fmt.Println("======>>> init success")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		_appInst.SetError(err.Error())
		return false
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	cert, err := loadLocalTlsConf(baseDir)
	if err != nil {
		_appInst.SetError(err.Error())
		return false
	}
	_appInst.localTls = &tls.Config{Certificates: []tls.Certificate{cert}}
	return true
}

func loadLocalTlsConf(baseDir string) (tls.Certificate, error) {
	certFile := path.Join(baseDir, "key.pem")
	keyFile := path.Join(baseDir, "cert.pem")
	_, ok1 := utils.FileExists(certFile)
	_, ok2 := utils.FileExists(keyFile)
	if ok2 && ok1 {
		return tls.LoadX509KeyPair(certFile, keyFile)
	}
	if err := utils.GenerateByParam(certFile, keyFile, 365, "", ""); err != nil {
		return tls.Certificate{}, err
	}

	return tls.LoadX509KeyPair(certFile, keyFile)
}

//export OpenWallet
func OpenWallet(auth string) bool {
	return false
}

//export ShowBalance
func ShowBalance(addr string) *C.char {
	return nil
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

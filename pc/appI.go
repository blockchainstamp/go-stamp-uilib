package main

/*
#include "callback.h"
*/
import "C"
import (
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bStamp "github.com/blockchainstamp/go-stamp-wallet"
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
		fmt.Println("======>>> int bStamp sdk err:", err)
		_appInst.SetError(err.Error())
		return false
	}
	fmt.Println("======>>> init bStamp SDK success")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Println("======>>> set log err:", err)
		_appInst.SetError(err.Error())
		return false
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	fmt.Println("======>>> set log level:", level)

	cert, err := loadLocalTlsConf(baseDir)
	if err != nil {
		fmt.Println("======>>> load tsl conf err:", err)
		_appInst.SetError(err.Error())
		return false
	}
	fmt.Println("======>>> tls config load success:")
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
		fmt.Println("======>>> create tls file err:", err)
		return tls.Certificate{}, err
	}
	fmt.Println("======>>> create tls files:", certFile, keyFile)

	return tls.LoadX509KeyPair(certFile, keyFile)
}

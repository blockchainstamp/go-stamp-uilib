package main

import "C"
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/protocol/imap"
	"github.com/blockchainstamp/go-mail-proxy/protocol/smtp"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel          string
	AllowInsecureAuth bool
	SMTPConf          *smtp.Conf
	IMAPConf          *imap.Conf
	CmdSrvAddr        string
}
type Service struct {
	imapSrv   *imap.Service
	smtpSrv   *smtp.Service
	srvStatus bool
	sigCh     chan struct{}
}

//export ServiceStatus
func ServiceStatus() bool {
	return _appInst.service != nil && _appInst.service.srvStatus
}

//export InitService
func InitService(confData []byte, auth string) bool {
	conf := &Config{}
	if err := json.Unmarshal(confData, conf); err != nil {
		_appInst.SetError(err.Error())
		logrus.Error(err.Error())
		return false
	}

	level, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		_appInst.SetError(err.Error())
		logrus.Error(err.Error())
		return false
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(_appInst.logger)

	var localTlsCfg *tls.Config
	if !conf.AllowInsecureAuth {
		logrus.Info("need local tls config")
		localTlsCfg = _appInst.localTls
	}

	imapSrv, err := imap.NewIMAPSrv(conf.IMAPConf, localTlsCfg)
	if err != nil {
		_appInst.SetError(err.Error())
		logrus.Error(err.Error())
		return false
	}

	w, err := bstamp.Inst().ActiveWallet(comm.WalletAddr(conf.SMTPConf.StampWalletAddr), auth)
	if err != nil {
		_appInst.SetError(err.Error())
		logrus.Error(err.Error())
		return false
	}
	logrus.Info("wallet address:", w.Address())
	logrus.Info("eth address:", w.EthAddr())

	smtpSrv, err := smtp.NewSMTPSrv(conf.SMTPConf, localTlsCfg)
	if err != nil {
		_appInst.SetError(err.Error())
		logrus.Error(err.Error())
		return false
	}

	_appInst.service = &Service{
		imapSrv:   imapSrv,
		smtpSrv:   smtpSrv,
		srvStatus: false,
		sigCh:     make(chan struct{}, 2),
	}
	go utils.StartCmdService(conf.CmdSrvAddr)

	_appInst.cfg = conf

	logrus.Info("init system success")
	return true
}

//export StartService
func StartService() bool {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[StartService()]=>service instance is nil")
		return false
	}
	if err := startSrv(); err != nil {
		_appInst.SetError(err.Error())
		logrus.Error(err.Error())
		return false
	}
	go monitor()
	return false
}
func monitor() {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[monitor()]=>service instance is nil")
		return
	}
	for {
		select {
		case <-srv.sigCh:
			srv.srvStatus = false
			msg := &CallBackMsg{
				Cmd:   CMDSrvStatusChanged,
				Param: fmt.Sprintf("%t", false),
			}
			_appInst.CBJsonData(msg)
			logrus.Warn("service stop by signal")
		}
	}
}

//export StopService
func StopService() {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[StopService()]=>service instance is nil")
		return
	}
	_appInst.service = nil
}

func startSrv() error {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[startSrv()]=>service instance is nil")
		return fmt.Errorf("service instance is nil")
	}
	var err error = nil
	if err = srv.imapSrv.Start(srv.sigCh); err != nil {
		return err
	}
	if err = srv.smtpSrv.Start(srv.sigCh); err != nil {
		return err
	}
	return err
}
func shutdown() {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[shutdown()]=>service instance is nil")
		return
	}
	srv.sigCh <- struct{}{}
}

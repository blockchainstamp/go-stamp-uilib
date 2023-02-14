package main

import "C"
import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/protocol/imap"
	"github.com/blockchainstamp/go-mail-proxy/protocol/smtp"
	"github.com/blockchainstamp/go-mail-proxy/utils"
	bstamp "github.com/blockchainstamp/go-stamp-wallet"
	"github.com/blockchainstamp/go-stamp-wallet/comm"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

type Config struct {
	LogLevel          string     `json:"log_level"`
	AllowInsecureAuth bool       `json:"allow_insecure_auth"`
	SMTPConf          *smtp.Conf `json:"smtp_conf"`
	IMAPConf          *imap.Conf `json:"imap_conf"`
	CmdSrvAddr        string     `json:"cmd_srv_addr"`
}

func (c *Config) String() string {
	s := "\n+++++++++++++++++++++++config+++++++++++++++++++++++++++++"
	s += "\nLog Level:\t" + c.LogLevel
	s += "\nSMTP Config:\t" + c.SMTPConf.String()
	s += "\nIMAP Config:\t" + c.IMAPConf.String()
	s += "\nCMD Srv Addr:\t" + c.CmdSrvAddr
	s += fmt.Sprintf("\nSecure Auth:\t%t", c.AllowInsecureAuth)
	s += "\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n"
	return s
}

type Service struct {
	imapSrv *imap.Service
	smtpSrv *smtp.Service
	ctx     context.Context
	cancel  context.CancelFunc
}

//export WriteCaFile
func WriteCaFile(name, data string) *C.char {
	caFile := path.Join(_appInst.basDir, name+".cer")
	_, ok := utils.FileExists(caFile)
	if ok {
		fmt.Println("file exist", caFile)
	}
	if err := os.WriteFile(caFile, []byte(data), 0600); err != nil {
		_appInst.SetError(err.Error())
		return nil
	}

	return C.CString(caFile)
}

//export ServiceStatus
func ServiceStatus() bool {
	return _appInst.service != nil
}

//export InitService
func InitService(confJson string, auth string) bool {
	conf := &Config{}
	if err := json.Unmarshal([]byte(confJson), conf); err != nil {
		_appInst.SetError(err.Error())
		return false
	}
	fmt.Println("----------------------------")
	fmt.Println(conf.String())
	fmt.Println("----------------------------")

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
		logrus.Info("prepare local tls config")
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
	ctx, cancel := context.WithCancel(context.Background())
	_appInst.service = &Service{
		imapSrv: imapSrv,
		smtpSrv: smtpSrv,
		ctx:     ctx,
		cancel:  cancel,
	}
	go utils.StartCmdService(conf.CmdSrvAddr, ctx)

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
	return true
}

func monitor() {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[monitor()]=>service instance is nil")
		return
	}
	for {
		select {
		case <-srv.ctx.Done():
			StopService()
			logrus.Warn("service stop by context")
			return
		}
	}
}

//export StopService
func StopService() {
	var srv = _appInst.service
	if srv == nil {
		logrus.Info("[StopService()]=>service instance is nil")
		return
	}
	_appInst.service.cancel()
	shutdown()
	_appInst.service = nil
	msg := &CallBackMsg{
		Cmd:   CMDSrvStatusChanged,
		Param: fmt.Sprintf("%t", false),
	}
	_appInst.CBJsonData(msg)
}

func startSrv() error {
	var srv = _appInst.service
	if srv == nil {
		logrus.Error("[startSrv()]=>service instance is nil")
		return fmt.Errorf("service instance is nil")
	}
	var err error = nil
	if err = srv.imapSrv.StartWithCtx(_appInst.service.cancel); err != nil {
		return err
	}
	if err = srv.smtpSrv.StartWithCtx(_appInst.service.cancel); err != nil {
		return err
	}
	return err

}

func shutdown() {
	var srv = _appInst.service
	if srv == nil {
		return
	}
	_appInst.service.smtpSrv.Stop()
	_appInst.service.imapSrv.Stop()
}

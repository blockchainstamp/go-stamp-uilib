package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/blockchainstamp/go-mail-proxy/protocol/common"
	contract "github.com/blockchainstamp/go-stamp-wallet/smart_contract"
	"path"
	"strings"
	"testing"
)

var (
	addr = ""
)

func init() {
	flag.StringVar(&addr, "addr", "0xF9Cbfd74808f812a3B8A2337BFC426B2A10Bd74a", "--addr")
}
func TestStampBalanceOfWallet(t *testing.T) {
	fmt.Println(path.Clean("file:///Users/hyperorchid/Library/Containers/com.simpleorg.bstamp.BStamp/Data/Documents/ribencong@126.com.cer"))
	c := &Config{
		LogLevel:   "info",
		CmdSrvAddr: "file:///Users/hyperorchid/Library/Containers/com.simpleorg.bstamp.BStamp/Data/Documents/ribencong@126.com.cer",
	}
	bts, _ := json.Marshal(c)
	print(string(bts))

	ss := `
{"cmd_srv_addr":"127.0.0.1:2100","smtp":{"read_time_out":10,"srv_addr":":443","max_recipients":50,"remote_conf":{"99927800@qq.com":{"ca_files":"","ca_domain":"smtp.qq.com","remote_srv_name":"smtp.qq.com","active_stamp_addr":"0xF9Cbfd74808f812a3B8A2337BFC426B2A10Bd74a","allow_not_secure":true,"remote_srv_port":443},"ribencong@163.com":{"active_stamp_addr":"","ca_domain":"smtp.163.com","ca_files":"file:///Users/hyperorchid/Library/Containers/com.simpleorg.bstamp.BStamp/Data/Documents/ribencong@163.com.cer","allow_not_secure":true,"remote_srv_port":443,"remote_srv_name":"smtp.163.com"},"ribencong@126.com":{"remote_srv_port":443,"ca_domain":"smtp.126.com","allow_not_secure":true,"remote_srv_name":"smtp.126.com","active_stamp_addr":"0xF9Cbfd74808f812a3B8A2337BFC426B2A10Bd74a","ca_files":"file:///Users/hyperorchid/Library/Containers/com.simpleorg.bstamp.BStamp/Data/Documents/ribencong@126.com.cer"},"ribencong@foxmail.com":{"active_stamp_addr":"","remote_srv_port":443,"ca_domain":"smtp.foxmail.com","remote_srv_name":"smtp.foxmail.com","allow_not_secure":false,"ca_files":""}},"max_message_bytes":1073741824,"write_time_out":10,"stamp_wallet_addr":"BS9G93bX1zMoN5A5ZQBwUf","SrvDomain":"localhost"},"log_level":"info","imap_conf":{"remote_conf":{"ribencong@126.com":{"remote_srv_name":"imap.126.com","remote_srv_port":996,"ca_domain":"imap.126.com","allow_not_secure":true,"ca_files":"file:///Users/hyperorchid/Library/Containers/com.simpleorg.bstamp.BStamp/Data/Documents/ribencong@126.com.cer"},"ribencong@163.com":{"allow_not_secure":true,"ca_files":"file:///Users/hyperorchid/Library/Containers/com.simpleorg.bstamp.BStamp/Data/Documents/ribencong@163.com.cer","ca_domain":"imap.163.com","remote_srv_name":"imap.163.com","remote_srv_port":996},"ribencong@foxmail.com":{"ca_domain":"imap.foxmail.com","ca_files":"","remote_srv_name":"imap.foxmail.com","allow_not_secure":false,"remote_srv_port":996},"99927800@qq.com":{"remote_srv_name":"imap.qq.com","allow_not_secure":false,"ca_domain":"imap.qq.com","ca_files":"","remote_srv_port":996}},"srv_addr":":443","srv_domain":"localhost"},"allow_insecure_auth":true}
`
	if e := json.Unmarshal([]byte(ss), c); e != nil {
		t.Fatal(e)
	}
	print(c.IMAPConf.String())

	fileNames := strings.Split(c.IMAPConf.RemoteConf["ribencong@163.com"].RemoteSrvCAs, common.CAFileSep)
	if len(fileNames) == 0 {
		t.Fatal("can't read")
	}
	fmt.Println(fileNames[0])
}

func TestStampConfig(t *testing.T) {
	fmt.Println(contract.StampConfigFromBlockChain(addr))
}

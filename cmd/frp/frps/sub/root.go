package frpslib

import (
	//"fmt"
	"github.com/fatedier/golib/crypto"

	"github.com/fatedier/frp/cmd/frps/sub"
	_ "github.com/fatedier/frp/pkg/metrics"
	_ "github.com/fatedier/frp/web/frps"
)

func Run(cfgFilePath string) {
	crypto.DefaultSalt = "frp"

	sub.RunServerDefault(cfgFilePath)
}

// RunDir 在指定目录下按多个配置文件启动多个 frps 实例。
// 语义与 frpc 侧的 RunDir 类似：遍历目录中所有文件，每个文件启动一个服务。
func RunDir(cfgDir string) {
	crypto.DefaultSalt = "frp"
	_ = sub.RunMultipleServersDefault(cfgDir)
}

func Close(uid string) (ret bool) {
	return sub.SvrClose(uid)
}

func GetUids() (ret string) {
	return sub.GetSvrUids()
}

func IsRunning(uid string) (ret bool) {
	return sub.SvrIsRunning(uid)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RunContent(uid string, cfgContent string) (err string) {
	crypto.DefaultSalt = "frp"

	//fmt.Printf("RunContent uid %s\n", uid)

	err = ""
	s := sub.RunServerContent(cfgContent, uid)
	if s != nil {
		err = s.Error()
	}
	return
}

func RunFile(uid string, cfgFilePath string) (err string) {
	crypto.DefaultSalt = "frp"

	//fmt.Printf("RunContent uid %s\n", uid)

	err = ""
	s := sub.RunServerFile(cfgFilePath, uid)
	if s != nil {
		err = s.Error()
	}
	return
}

func RunContentWithForce(uid string, cfgContent string, forceRestart bool) (err string) {
	crypto.DefaultSalt = "frp"

	err = ""
	s := sub.RunServerContent(cfgContent, uid, forceRestart)
	if s != nil {
		err = s.Error()
	}
	return
}

func RunFileWithForce(uid string, cfgFilePath string, forceRestart bool) (err string) {
	crypto.DefaultSalt = "frp"

	err = ""
	s := sub.RunServerFile(cfgFilePath, uid, forceRestart)
	if s != nil {
		err = s.Error()
	}
	return
}

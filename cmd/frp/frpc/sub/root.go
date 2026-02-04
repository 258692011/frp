package frpclib

import (
	//"fmt"
	_ "github.com/fatedier/frp/web/frpc"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/policy/security"
	"github.com/fatedier/golib/crypto"
)

func Run(cfgFilePath string, allowUnsafe ...string) {
	crypto.DefaultSalt = "frp"

	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	sub.RunClientDefault(cfgFilePath, unsafeFeatures)
}

func RunDir(cfgDir string, allowUnsafe ...string) {
	crypto.DefaultSalt = "frp"

	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	sub.RunMultipleClientsDefault(cfgDir, unsafeFeatures)
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

func RunContent(uid string, cfgContent string, allowUnsafe ...string) (err string) {
	crypto.DefaultSalt = "frp"

	//fmt.Printf("RunContent uid %s\n", uid)

	err = ""
	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	s := sub.RunClientContent(cfgContent, uid, unsafeFeatures)
	if s != nil {
		//fmt.Printf("FRPC Service Run Error %s\n", s)
		err = "FRPC Service Run Error"
	}
	return
}

func RunFile(uid string, cfgFilePath string, allowUnsafe ...string) (err string) {
	crypto.DefaultSalt = "frp"

	//fmt.Printf("RunContent uid %s\n", uid)

	err = ""
	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	s := sub.RunClientDefault(cfgFilePath, unsafeFeatures)
	if s != nil {
		//fmt.Printf("%s\n", s)
		err = "FRPC Service Run Error"
	}
	return
}

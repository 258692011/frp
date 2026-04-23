package frpclib

import (
	//"fmt"
	"github.com/fatedier/frp/cmd/frpc/sub"
	"github.com/fatedier/frp/pkg/policy/security"
	_ "github.com/fatedier/frp/web/frpc"
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
		err = s.Error()
	}
	return
}

func RunFile(uid string, cfgFilePath string, allowUnsafe ...string) (err string) {
	crypto.DefaultSalt = "frp"

	//fmt.Printf("RunContent uid %s\n", uid)

	err = ""
	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	s := sub.RunClientFile(cfgFilePath, uid, unsafeFeatures)
	if s != nil {
		err = s.Error()
	}
	return
}

func RunContentWithForce(uid string, cfgContent string, forceRestart bool, allowUnsafe ...string) (err string) {
	crypto.DefaultSalt = "frp"

	err = ""
	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	s := sub.RunClientContent(cfgContent, uid, unsafeFeatures, forceRestart)
	if s != nil {
		err = s.Error()
	}
	return
}

func RunFileWithForce(uid string, cfgFilePath string, forceRestart bool, allowUnsafe ...string) (err string) {
	crypto.DefaultSalt = "frp"

	err = ""
	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	s := sub.RunClientFile(cfgFilePath, uid, unsafeFeatures, forceRestart)
	if s != nil {
		err = s.Error()
	}
	return
}

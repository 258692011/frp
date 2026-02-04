package frpslib

import (
	//"fmt"
	"github.com/fatedier/golib/crypto"

	_ "github.com/fatedier/frp/assets/frps"
	"github.com/fatedier/frp/cmd/frps/sub"
	_ "github.com/fatedier/frp/pkg/metrics"
)

func Run(cfgFilePath string) {
	crypto.DefaultSalt = "frp"

	sub.RunServerDefault(cfgFilePath)
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
		//fmt.Printf("FRPS Service Run Error %s\n", s)
		err = "FRPS Service Run Error"
	}
	return
}

func RunFile(uid string, cfgFilePath string) (err string) {
	crypto.DefaultSalt = "frp"

	//fmt.Printf("RunContent uid %s\n", uid)

	err = ""
	s := sub.RunServerFile(cfgFilePath, uid)
	if s != nil {
		//fmt.Printf("%s\n", s)
		err = "FRPS Service Run Error"
	}
	return
}

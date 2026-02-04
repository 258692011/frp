package client

import (
	frpclib "github.com/fatedier/frp/cmd/frp/frpc/sub"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(5)
}

func Run(cfgFilePath string) {
	frpclib.Run(cfgFilePath)
}

func RunDir(cfgDir string) {
	frpclib.RunDir(cfgDir)
}

func Close(uid string) (ret bool) {
	defer debug.FreeOSMemory()
	ret = frpclib.Close(uid)
	return
}

func GetUids() (ret string) {
	ret = frpclib.GetUids()
	return
}

func IsRunning(uid string) (ret bool) {
	ret = frpclib.IsRunning(uid)
	return
}

func RunContent(uid string, cfgContent string) (err string) {
	err = frpclib.RunContent(uid, cfgContent)
	return
}

func RunFile(uid string, cfgFilePath string) (err string) {
	err = frpclib.RunFile(uid, cfgFilePath)
	return
}

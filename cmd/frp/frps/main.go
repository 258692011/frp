package server

import (
	frpslib "github.com/fatedier/frp/cmd/frp/frps/sub"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(5)
}

func Run(cfgFilePath string) {
	frpslib.Run(cfgFilePath)
}

func RunDir(cfgDir string) {
	frpslib.RunDir(cfgDir)
}

func Close(uid string) (ret bool) {
	defer debug.FreeOSMemory()
	ret = frpslib.Close(uid)
	return
}

func GetUids() (ret string) {
	ret = frpslib.GetUids()
	return
}

func IsRunning(uid string) (ret bool) {
	ret = frpslib.IsRunning(uid)
	return
}

func RunContent(uid string, cfgContent string) (err string) {
	err = frpslib.RunContent(uid, cfgContent)
	return
}

func RunFile(uid string, cfgFilePath string) (err string) {
	err = frpslib.RunFile(uid, cfgFilePath)
	return
}

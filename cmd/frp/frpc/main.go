package client

import (
	frpclib "github.com/fatedier/frp/cmd/frp/frpc/sub"
	"runtime/debug"

	// gobind / gomobile bind 需要该包出现在模块依赖图中（go mod tidy 不会从间接依赖推断出 bind）。
	_ "golang.org/x/mobile/bind"
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

func RunContentWithForce(uid string, cfgContent string, forceRestart bool) (err string) {
	err = frpclib.RunContentWithForce(uid, cfgContent, forceRestart)
	return
}

func RunFileWithForce(uid string, cfgFilePath string, forceRestart bool) (err string) {
	err = frpclib.RunFileWithForce(uid, cfgFilePath, forceRestart)
	return
}

// for android

package sub

import (
	"context"
	"fmt"
	//"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/policy/security"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/server"
)

// Local copies of the main package flags/configs for Android usage.
// We keep them independent from cmd/frps so gomobile can bind this
// package without depending on the main package globals.
var (
	cfgFile          string
	strictConfigMode = true
	allowUnsafe      []string

	serverCfg v1.ServerConfig

	Svrs = make(map[string]*server.Service)
)

func GetSvrUids() (ret string) {
	v := os.Getuid()
	v = v
	//fmt.Printf("os.Getuid() is %+v\n", v)
	if len(Svrs) > 0 {
		var a string = ""
		for k, _ := range Svrs {
			if a != "" {
				a += ","
			}
			a += k
		}
		ret = a
	} else {
		ret = ""
	}
	//fmt.Printf("GetSvrUids %+v\n", ret)
	return
}

func SvrIsRunning(uid string) (ret bool) {
	//fmt.Printf("︷︷︷︷︷︷︷︷︷︷︷\n")
	//fmt.Printf("SvrIsRunning uid %+v\n", uid)
	s := false
	v, ok := Svrs[uid]
	v = v
	//fmt.Printf("SvrIsRunning Svrs[uid] %+v is %+v\n", uid, v)
	if ok && Svrs[uid] != nil {
		s = true
	} else {
		SvrClose(uid)
	}
	//fmt.Printf("SvrIsRunning uid %+v is %+v\n", uid, s)
	GetSvrUids()
	//fmt.Printf("︸︸︸︸︸︸︸︸︸︸︸\n")
	ret = s
	return
}

func SvrClose(uid string) (ret bool) {
	s := true
	if svr, ok := Svrs[uid]; ok {
		if svr != nil {
			_ = svr.Close()
			//fmt.Printf("SvrClose Svrs[uid] %+v\n", uid)
		}
		delete(Svrs, uid)
		//fmt.Printf("delete Svrs uid %+v\n", uid)
	}
	GetSvrUids()
	ret = s
	return
}

func RunServerDefault(cfgFilePath string) (err error) {
	cfgFile = cfgFilePath
	var (
		svrCfg         *v1.ServerConfig
		isLegacyFormat bool
	)
	if cfgFile != "" {
		svrCfg, isLegacyFormat, err = config.LoadServerConfig(cfgFile, strictConfigMode)
		if err != nil {
			//fmt.Println(err)
			os.Exit(1)
		}
		if isLegacyFormat {
			//fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " + "please use yaml/json/toml format instead!\n")
		}
	} else {
		serverCfg.Complete()
		svrCfg = &serverCfg
	}

	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	validator := validation.NewConfigValidator(unsafeFeatures)
	warning, err := validator.ValidateServerConfig(svrCfg)
	if warning != nil {
		//fmt.Printf("WARNING: %+v\n", warning)
	}
	if err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
	cfg := svrCfg
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	if cfgFile != "" {
		log.Infof("frps uses config file: %+v", cfgFile)
	} else {
		log.Infof("frps uses command line arguments for config")
	}

	svr, err := server.NewService(cfg)
	if err != nil {
		return err
	}
	log.Infof("frps started successfully")
	svr.Run(context.Background())
	return
}

// RunMultipleServersDefault 启动指定目录下的多份 frps 配置文件。
// 语义与 frpc 侧的 RunMultipleClientsDefault 类似：
// - 遍历 cfgDir 中的所有普通文件（忽略子目录）
// - 每个文件作为一份独立的 frps 配置启动一个 Service 实例
// - 所有启动 goroutine 结束后函数返回
func RunMultipleServersDefault(cfgDir string) error {
	var wg sync.WaitGroup

	err := filepath.WalkDir(cfgDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		wg.Add(1)
		// 略微错开放启动时间，避免同时大量起服
		time.Sleep(time.Millisecond)

		go func(cfgPath string) {
			defer wg.Done()
			// 使用配置文件路径作为 uid 方便在 Svrs 里区分，这里忽略具体错误
			_ = RunServerFile(cfgPath, cfgPath)
		}(path)

		return nil
	})

	wg.Wait()
	return err
}

func tempFile(content string) (path string) { //android下会无权限，就算给app文件权限也没用
	// 创建临时文件
	file, err := ioutil.TempFile("", "temp") // 第二个参数为前缀名称
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 将内容写入临时文件
	_, err = file.Write([]byte(content)) //_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}

	// 获取临时文件路径
	path = file.Name()
	//fmt.Println("frpc temp file: ", path)
	return
}

func shouldForceRestart(forceRestart []bool) bool {
	return len(forceRestart) > 0 && forceRestart[0]
}

func ensureServerSlot(uid string, forceRestart bool) error {
	if uid == "" {
		return fmt.Errorf("uid cannot be empty")
	}
	if existing, ok := Svrs[uid]; ok && existing != nil {
		if !forceRestart {
			return fmt.Errorf("service already running for uid: %s", uid)
		}
		SvrClose(uid)
	}
	return nil
}

func RunServerContent(cfgContent string, uid string, forceRestart ...bool) (err error) {
	if err := ensureServerSlot(uid, shouldForceRestart(forceRestart)); err != nil {
		return err
	}

	var (
		svrCfg         *v1.ServerConfig
		isLegacyFormat bool
	)
	svrCfg, isLegacyFormat, err = config.LoadServerConfigFromContent(cfgContent, strictConfigMode)
	if err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
	if isLegacyFormat {
		//fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " + "please use yaml/json/toml format instead!\n")
	}

	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	validator := validation.NewConfigValidator(unsafeFeatures)
	warning, err := validator.ValidateServerConfig(svrCfg)
	if warning != nil {
		//fmt.Printf("WARNING: %+v\n", warning)
	}
	if err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
	cfg := svrCfg
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	log.Infof("frps uses content for config")

	svr, err := server.NewService(cfg)
	if err != nil {
		return err
	}
	log.Infof("frps started successfully")
	Svrs[uid] = svr
	go func() {
		svr.Run(context.Background())
		SvrClose(uid)
	}()
	return
}

func RunServerFile(cfgFilePath string, uid string, forceRestart ...bool) (err error) {
	if err := ensureServerSlot(uid, shouldForceRestart(forceRestart)); err != nil {
		return err
	}

	cfgFile = cfgFilePath
	var (
		svrCfg         *v1.ServerConfig
		isLegacyFormat bool
	)
	if cfgFile != "" {
		svrCfg, isLegacyFormat, err = config.LoadServerConfig(cfgFile, strictConfigMode)
		if err != nil {
			//fmt.Println(err)
			os.Exit(1)
		}
		if isLegacyFormat {
			//fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " + "please use yaml/json/toml format instead!\n")
		}
	} else {
		serverCfg.Complete()
		svrCfg = &serverCfg
	}

	unsafeFeatures := security.NewUnsafeFeatures(allowUnsafe)
	validator := validation.NewConfigValidator(unsafeFeatures)
	warning, err := validator.ValidateServerConfig(svrCfg)
	if warning != nil {
		//fmt.Printf("WARNING: %+v\n", warning)
	}
	if err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
	cfg := svrCfg
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	if cfgFile != "" {
		log.Infof("frps uses config file: %+v", cfgFile)
	} else {
		log.Infof("frps uses command line arguments for config")
	}

	svr, err := server.NewService(cfg)
	if err != nil {
		return err
	}
	log.Infof("frps started successfully")
	Svrs[uid] = svr
	go func() {
		svr.Run(context.Background())
		SvrClose(uid)
	}()
	return
}

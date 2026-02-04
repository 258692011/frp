// for android

package sub

import (
	"context"

	//"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/policy/featuregate"
	"github.com/fatedier/frp/pkg/policy/security"
	"github.com/fatedier/frp/pkg/util/log"
)

var (
	svrs = make(map[string]*client.Service)
)

func GetSvrUids() (ret string) {
	v := os.Getuid()
	v = v
	//fmt.Printf("os.Getuid() is %+v\n", v)
	if len(svrs) > 0 {
		var a string = ""
		for k := range svrs {
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
	v, ok := svrs[uid]
	v = v
	//fmt.Printf("SvrIsRunning svrs[uid] %+v is %+v\n", uid, v)
	if ok && svrs[uid] != nil {
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
	if svr, ok := svrs[uid]; ok {
		if svr != nil {
			svr.Close()
			//fmt.Printf("SvrClose svrs[uid] %+v\n", uid)
		}
		delete(svrs, uid)
		//fmt.Printf("delete svrs uid %+v\n", uid)
	}
	GetSvrUids()
	ret = s
	return
}

func isProcessExist(pid int) bool {
	cmd := exec.Command("ps", "aux")
	out, err := cmd.Output()
	if err != nil {
		//fmt.Println(err)
		return false
	}
	output := string(out)
	processes := strings.Split(output, "\n")
	for _, process := range processes {
		fields := strings.Fields(process)
		if len(fields) > 1 {
			processPid, err := strconv.Atoi(fields[1])
			if err == nil && processPid == pid {
				return true
			}
		}
	}
	return false
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

func RunMultipleClientsDefault(cfgDir string, unsafeFeatures *security.UnsafeFeatures) error {
	var wg sync.WaitGroup
	err := filepath.WalkDir(cfgDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		wg.Add(1)
		time.Sleep(time.Millisecond)
		go func() {
			defer wg.Done()
			err := runClient(path, unsafeFeatures)
			if err != nil {
				//fmt.Printf("frpc service error for config file [%s]\n", path)
			}
		}()
		return nil
	})
	wg.Wait()
	return err
}

func RunClientDefault(cfgFilePath string, unsafeFeatures *security.UnsafeFeatures) error {
	return runClient(cfgFilePath, unsafeFeatures)
}

func RunClientContent(cfgContent string, uid string, unsafeFeatures *security.UnsafeFeatures) error {
	cfg, proxyCfgs, visitorCfgs, isLegacyFormat, err := config.LoadClientConfigFromContent(cfgContent, strictConfigMode)
	if err != nil {
		return err
	}
	if isLegacyFormat {
		//fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " + "please use yaml/json/toml format instead!\n")
	}

	if len(cfg.FeatureGates) > 0 {
		if err := featuregate.SetFromMap(cfg.FeatureGates); err != nil {
			return err
		}
	}

	warning, err := validation.ValidateAllClientConfig(cfg, proxyCfgs, visitorCfgs, unsafeFeatures)
	if warning != nil {
		//fmt.Printf("WARNING: %+v\n", warning)
	}
	if err != nil {
		return err
	}
	err = startServiceContent(cfg, proxyCfgs, visitorCfgs, unsafeFeatures, "", uid)
	if err != nil {
		//fmt.Printf("%+v\n", err)
		SvrClose(uid)
	}
	return err
}

func handleTermSignalWithID(svr *client.Service, doneCh chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
	svr = nil
	close(doneCh)
}

func startServiceContent(
	cfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
	unsafeFeatures *security.UnsafeFeatures,
	cfgFile string,
	uid string,
) error {
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	if cfgFile != "" {
		log.Infof("start frpc service for config file [%+v]", cfgFile)
		defer log.Infof("frpc service for config file [%+v] stopped", cfgFile)
	}

	svr, err := client.NewService(client.ServiceOptions{
		Common:         cfg,
		ProxyCfgs:      proxyCfgs,
		VisitorCfgs:    visitorCfgs,
		UnsafeFeatures: unsafeFeatures,
		ConfigFilePath: cfgFile,
	})
	if err != nil {
		return err
	}

	svrs[uid] = svr

	closedDoneCh := make(chan struct{})
	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	// Capture the exit signal if we use kcp or quic.
	if shouldGracefulClose {
		go handleTermSignalWithID(svrs[uid], closedDoneCh)
	}

	err = svrs[uid].Run(context.Background())
	if err == nil && shouldGracefulClose {
		<-closedDoneCh
	}
	return err
}


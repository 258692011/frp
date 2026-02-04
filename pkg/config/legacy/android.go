// for android

package legacy

import (
	"bytes"
	"fmt"
)

func ParseClientConfigFromContent(cfgContent string) (
	cfg ClientCommonConf,
	proxyCfgs map[string]ProxyConf,
	visitorCfgs map[string]VisitorConf,
	err error,
) {
	var content []byte
	content = []byte(cfgContent)
	if err != nil {
		return
	}
	configBuffer := bytes.NewBuffer(nil)
	configBuffer.Write(content)

	// Parse common section.
	cfg, err = UnmarshalClientConfFromIni(content)
	if err != nil {
		return
	}
	if err = cfg.Validate(); err != nil {
		err = fmt.Errorf("parse config error: %v", err)
		return
	}

	// Aggregate proxy configs from include files.
	var buf []byte
	buf, err = getIncludeContents(cfg.IncludeConfigFiles)
	if err != nil {
		err = fmt.Errorf("getIncludeContents error: %v", err)
		return
	}
	configBuffer.WriteString("\n")
	configBuffer.Write(buf)

	// Parse all proxy and visitor configs.
	proxyCfgs, visitorCfgs, err = LoadAllProxyConfsFromIni(cfg.User, configBuffer.Bytes(), cfg.Start)
	if err != nil {
		return
	}
	return
}


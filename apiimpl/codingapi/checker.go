package codingapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/arthurkiller/rollingwriter"
	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-gateway/utils"
	"github.com/traefik/yaegi/interp"
)

type CheckScript struct {
	engineChan chan *interp.Interpreter
}

func newCheckScript(path string) (*CheckScript, error) {
	engineChan := make(chan *interp.Interpreter, 16)
	for i := 0; i < 16; i++ {
		script, err := utils.LoadScript(path)
		if err != nil {
			close(engineChan)
			return nil, err
		}
		engineChan <- script
	}
	return &CheckScript{
		engineChan: engineChan,
	}, nil
}

func (script *CheckScript) Exec(apiUrl, content string) (bool, error) {
	engine := <-script.engineChan
	defer func() {
		script.engineChan <- engine
	}()

	f, err := engine.Eval("script.CheckContent")
	if err != nil {
		return false, err
	}
	retList := f.Call([]reflect.Value{reflect.ValueOf(apiUrl), reflect.ValueOf(content)})
	if len(retList) != 1 {
		return false, fmt.Errorf("wrong return type in check content")
	}
	return retList[0].Bool(), nil
}

type ContentChecker struct {
	globalChecker   *CheckScript
	codingCheckList []*CheckScript
	logWriter       rollingwriter.RollingWriter
}

func newContentChecker(cfg *config.GateWayConfig, logWriter rollingwriter.RollingWriter) (*ContentChecker, error) {
	var globalChecker *CheckScript
	codingCheckList := []*CheckScript{}
	var err error

	if cfg.CheckScript != "" {
		globalChecker, err = newCheckScript(cfg.CheckScript)
		if err != nil {
			return nil, err
		}
	}
	for _, codingProvider := range cfg.Provider.CodingProviderList {
		if codingProvider.Checkscript == "" {
			codingCheckList = append(codingCheckList, nil)
		} else {
			checkScript, err := newCheckScript(codingProvider.Checkscript)
			if err != nil {
				return nil, err
			}
			codingCheckList = append(codingCheckList, checkScript)
		}
	}

	return &ContentChecker{
		globalChecker:   globalChecker,
		codingCheckList: codingCheckList,
		logWriter:       logWriter,
	}, nil
}

func (checker *ContentChecker) CheckCodeContent(apiUrl string, providerIndex int, content string) (bool, error) {
	var checkScript *CheckScript
	if providerIndex >= len(checker.codingCheckList) {
		checkScript = checker.globalChecker
	} else {
		checkScript = checker.codingCheckList[providerIndex]
		if checkScript == nil {
			checkScript = checker.globalChecker
		}
	}
	//默认通过
	if checkScript == nil {
		return true, nil
	}
	allow, err := checkScript.Exec(apiUrl, content)
	if err != nil {
		return allow, err
	}
	obj := map[string]interface{}{
		"apiUrl":  apiUrl,
		"content": content,
		"allow":   allow,
	}
	logData, err := json.Marshal(obj)
	if err == nil {
		timeStr := time.Now().Format(time.RFC3339)
		fmt.Fprintf(checker.logWriter, "%s %s", timeStr, string(logData))
	}
	return allow, nil
}

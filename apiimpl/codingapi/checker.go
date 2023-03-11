package codingapi

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-gateway/utils"
	"github.com/traefik/yaegi/interp"
)

type CheckScript struct {
	lock   sync.Mutex
	engine *interp.Interpreter
}

func newCheckScript(path string) (*CheckScript, error) {
	absPath, err := utils.GetAbsPath(path)
	if err != nil {
		return nil, err
	}
	engine := interp.New(interp.Options{})
	_, err = engine.EvalPath(absPath)
	if err != nil {
		return nil, err
	}
	return &CheckScript{
		engine: engine,
	}, nil
}

func (script *CheckScript) Exec(apiUrl, content string) (bool, error) {
	script.lock.Lock()
	defer script.lock.Unlock()

	f, err := script.engine.Eval("script.CheckContent")
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
}

func newContentChecker(cfg *config.GateWayConfig) (*ContentChecker, error) {
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
	return checkScript.Exec(apiUrl, content)
}

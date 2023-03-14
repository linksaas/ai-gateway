package codingapi

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-gateway/utils"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

type BackendScript struct {
	lock   sync.Mutex
	engine *interp.Interpreter
}

func newBackendScript(path string) (*BackendScript, error) {
	absPath, err := utils.GetAbsPath(path)
	if err != nil {
		return nil, err
	}
	engine := interp.New(interp.Options{
		GoPath:    os.Getenv("GOROOT"),
		BuildTags: []string{},
		Env:       os.Environ(),
	})
	err = engine.Use(stdlib.Symbols)
	if err != nil {
		return nil, err
	}
	err = engine.Use(interp.Symbols)
	if err != nil {
		return nil, err
	}
	err = engine.Use(syscall.Symbols)
	if err != nil {
		return nil, err
	}
	os.Setenv("YAEGI_SYSCALL", "1")
	err = engine.Use(unsafe.Symbols)
	if err != nil {
		return nil, err
	}
	os.Setenv("YAEGI_UNSAFE", "1")
	scriptData, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	_, err = engine.Eval(string(scriptData))
	if err != nil {
		return nil, err
	}
	return &BackendScript{
		engine: engine,
	}, nil
}

func (script *BackendScript) CallComplete(lang, content string) ([]string, error) {
	script.lock.Lock()
	defer script.lock.Unlock()

	f, err := script.engine.Eval("script.Complete")
	if err != nil {
		return nil, err
	}
	retList := f.Call([]reflect.Value{reflect.ValueOf(lang), reflect.ValueOf(content)})
	if len(retList) != 1 {
		return nil, fmt.Errorf("wrong return type in complete func")
	}
	contentListValue := retList[0]
	contentList := []string{}
	for i := 0; i < contentListValue.Len(); i++ {
		contentList = append(contentList, contentListValue.Index(i).String())
	}
	return contentList, nil
}

func (script *BackendScript) CallConvert(lang, destLang, content string) ([]string, error) {
	script.lock.Lock()
	defer script.lock.Unlock()

	f, err := script.engine.Eval("script.Convert")
	if err != nil {
		return nil, err
	}
	retList := f.Call([]reflect.Value{reflect.ValueOf(lang), reflect.ValueOf(destLang), reflect.ValueOf(content)})
	if len(retList) != 1 {
		return nil, fmt.Errorf("wrong return type in convert func")
	}
	contentListValue := retList[0]
	contentList := []string{}
	for i := 0; i < contentListValue.Len(); i++ {
		contentList = append(contentList, contentListValue.Index(i).String())
	}
	return contentList, nil
}

func (script *BackendScript) CallExplain(lang, content string) ([]string, error) {
	script.lock.Lock()
	defer script.lock.Unlock()

	f, err := script.engine.Eval("script.Explain")
	if err != nil {
		return nil, err
	}
	retList := f.Call([]reflect.Value{reflect.ValueOf(lang), reflect.ValueOf(content)})
	if len(retList) != 1 {
		return nil, fmt.Errorf("wrong return type in explain func")
	}
	contentListValue := retList[0]
	contentList := []string{}
	for i := 0; i < contentListValue.Len(); i++ {
		contentList = append(contentList, contentListValue.Index(i).String())
	}
	return contentList, nil
}

func (script *BackendScript) CallFixError(lang, errStr string) ([]string, error) {
	script.lock.Lock()
	defer script.lock.Unlock()

	f, err := script.engine.Eval("script.Fixerror")
	if err != nil {
		return nil, err
	}
	retList := f.Call([]reflect.Value{reflect.ValueOf(lang), reflect.ValueOf(errStr)})
	if len(retList) != 1 {
		return nil, fmt.Errorf("wrong return type in fixerror func")
	}
	contentListValue := retList[0]
	contentList := []string{}
	for i := 0; i < contentListValue.Len(); i++ {
		contentList = append(contentList, contentListValue.Index(i).String())
	}
	return contentList, nil
}

func (script *BackendScript) CallGenTest(lang, content string) ([]string, error) {
	script.lock.Lock()
	defer script.lock.Unlock()

	f, err := script.engine.Eval("script.Gentest")
	if err != nil {
		return nil, err
	}
	retList := f.Call([]reflect.Value{reflect.ValueOf(lang), reflect.ValueOf(content)})
	if len(retList) != 1 {
		return nil, fmt.Errorf("wrong return type in gentest func")
	}
	contentListValue := retList[0]
	contentList := []string{}
	for i := 0; i < contentListValue.Len(); i++ {
		contentList = append(contentList, contentListValue.Index(i).String())
	}
	return contentList, nil
}

type ScriptBackend struct {
	scriptList []*BackendScript
}

func newScriptBackend(cfg *config.GateWayConfig) (*ScriptBackend, error) {
	scriptList := []*BackendScript{}
	for _, codingProvider := range cfg.Provider.CodingProviderList {
		if strings.HasPrefix(codingProvider.Backend, "script://") {
			path := strings.TrimPrefix(codingProvider.Backend, "script://")
			script, err := newBackendScript(path)
			if err != nil {
				return nil, err
			}
			scriptList = append(scriptList, script)
		} else {
			scriptList = append(scriptList, nil)
		}
	}
	return &ScriptBackend{scriptList: scriptList}, nil
}

func (backend *ScriptBackend) GetBackendScript(index int) *BackendScript {
	if index >= 0 && index < len(backend.scriptList) {
		return backend.scriptList[index]
	}
	return nil
}

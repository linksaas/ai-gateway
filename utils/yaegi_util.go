package utils

import (
	"os"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

func LoadScript(path string) (*interp.Interpreter, error) {
	absPath, err := GetAbsPath(path)
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
	return engine, nil
}

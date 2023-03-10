package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed config.yaml
var configTpl string

//go:embed script/check.go
var checkTpl string

//go:embed script/coding_provider.go
var codingProviderTpl string

var initCmd = &cobra.Command{
	Use:  "init",
	RunE: initCfgAndScript,
}

func initCfgAndScript(cmd *cobra.Command, args []string) error {
	_, err := os.Stat("config.yaml")
	if err == nil {
		return fmt.Errorf("config.yaml exist,please backup and remove it")
	}
	err = ioutil.WriteFile("config.yaml", []byte(configTpl), 0600)
	if err != nil {
		return err
	}

	os.MkdirAll("script", 0700) //skip check error

	checkScriptPath := strings.Join([]string{"script", "check.go"}, string(os.PathSeparator))
	_, err = os.Stat(checkScriptPath)
	if err == nil {
		return fmt.Errorf("%s exist,please backup and remove it", checkScriptPath)
	}
	err = ioutil.WriteFile(checkScriptPath, []byte(checkTpl), 0600)
	if err != nil {
		return err
	}

	codingProviderPath := strings.Join([]string{"script", "coding_provider.go"}, string(os.PathSeparator))
	_, err = os.Stat(codingProviderPath)
	if err == nil {
		return fmt.Errorf("%s exist,please backup and remove it", codingProviderPath)
	}
	err = ioutil.WriteFile(codingProviderPath, []byte(codingProviderTpl), 0600)
	if err != nil {
		return err
	}
	return nil
}

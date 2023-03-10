package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/linksaas/ai-gateway/utils"
	"github.com/traefik/yaegi/interp"
	"gopkg.in/yaml.v3"
)

func tryLoadScript(path string) error {
	i := interp.New(interp.Options{})
	_, err := i.EvalPath(path)
	if err != nil {
		return err
	}
	return nil
}

type SslConfig struct {
	Enable bool   `yaml:"enable"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
}

func (cfg *SslConfig) check() error {
	if !cfg.Enable {
		return nil
	}
	certPath, err := utils.GetAbsPath(cfg.Cert)
	if err != nil {
		return err
	}
	_, err = os.Stat(certPath)
	if err != nil {
		return err
	}
	keyPath, err := utils.GetAbsPath(cfg.Key)
	if err != nil {
		return err
	}
	_, err = os.Stat(keyPath)
	if err != nil {
		return err
	}
	return nil
}

type CodingProviderConfig struct {
	Backend          string   `yaml:"backend"`
	Checkscript      string   `yaml:"checkscript"`
	CompleteLangList []string `yaml:"complete"`
	ConvertLangList  []string `yaml:"convert"`
	ExplainLangList  []string `yaml:"explain"`
	FixerrorLangList []string `yaml:"fixerror"`
	GentestLangList  []string `yaml:"gentest"`
}

func (cfg *CodingProviderConfig) checkLang(langList []string) error {
	validLangList := []string{
		"python",
		"c",
		"cplusplus",
		"java",
		"csharp",
		"visualbasic",
		"javascript",
		"sql",
		"asm",
		"php",
		"r",
		"go",
		"matlab",
		"swift",
		"delphi",
		"ruby",
		"perl",
		"objc",
		"rust",
	}
	for _, lang := range langList {
		lang = strings.ToLower(lang)
		match := false
		for _, validLang := range validLangList {
			if validLang == lang {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("unkwown lang %s", lang)
		}
	}
	return nil
}

func (cfg *CodingProviderConfig) check() error {
	if strings.HasPrefix(cfg.Backend, "http://") {
		//do nothing
	} else if strings.HasPrefix(cfg.Backend, "script://") {
		path, _ := strings.CutPrefix(cfg.Backend, "script://")
		absPath, err := utils.GetAbsPath(path)
		if err != nil {
			return err
		}
		err = tryLoadScript(absPath)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("only support http:// and script:// schema")
	}
	if cfg.Checkscript != "" {
		absPath, err := utils.GetAbsPath(cfg.Checkscript)
		if err != nil {
			return err
		}
		err = tryLoadScript(absPath)
		if err != nil {
			return err
		}
	}
	err := cfg.checkLang(cfg.CompleteLangList)
	if err != nil {
		return err
	}
	err = cfg.checkLang(cfg.ConvertLangList)
	if err != nil {
		return err
	}
	err = cfg.checkLang(cfg.ExplainLangList)
	if err != nil {
		return err
	}
	err = cfg.checkLang(cfg.FixerrorLangList)
	if err != nil {
		return err
	}
	err = cfg.checkLang(cfg.GentestLangList)
	if err != nil {
		return err
	}
	return nil
}

type ProviderConfig struct {
	CodingProviderList []CodingProviderConfig `yaml:"coding"`
}

func (cfg *ProviderConfig) check() error {
	for _, codingCfg := range cfg.CodingProviderList {
		err := codingCfg.check()
		if err != nil {
			return err
		}
	}
	return nil
}

type GateWayConfig struct {
	Port        uint16         `yaml:"port"`
	Ssl         SslConfig      `yaml:"ssl"`
	Secret      string         `yaml:"secret"`
	Dev         bool           `yaml:"dev"`
	TokenTtl    int            `yaml:"tokenttl"`
	CheckScript string         `yaml:"checkscript"`
	Provider    ProviderConfig `yaml:"provider"`
}

func (cfg *GateWayConfig) check() error {
	err := cfg.Ssl.check()
	if err != nil {
		return err
	}
	if cfg.CheckScript != "" {
		checkScriptPath, err := utils.GetAbsPath(cfg.CheckScript)
		if err != nil {
			return err
		}
		err = tryLoadScript(checkScriptPath)
		if err != nil {
			return err
		}
	}
	if !cfg.Dev && len(cfg.Secret) < 32 {
		return fmt.Errorf("secret must have 32 chars at least")
	}
	err = cfg.Provider.check()
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(fileName string) (*GateWayConfig, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfgData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg := &GateWayConfig{}
	err = yaml.Unmarshal(cfgData, cfg)
	if err != nil {
		return nil, err
	}
	err = cfg.check()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

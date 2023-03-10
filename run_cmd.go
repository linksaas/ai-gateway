package main

import (
	"embed"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/apiimpl/codingapi"
	"github.com/linksaas/ai-gateway/apiimpl/devapi"
	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-gateway/utils"
	"github.com/spf13/cobra"
)

var runConfigFile string

//go:embed web/dist
var webFs embed.FS

var runCmd = &cobra.Command{
	Use:  "run",
	RunE: runGateWay,
}

func initWeb(engine *gin.Engine) {
	engine.GET("/*filepath", func(c *gin.Context) {
		fileName := c.Param("filepath")
		if fileName == "/" {
			fileName = "/index.html"
		}

		f, err := webFs.Open("web/dist" + fileName)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		ctype := mime.TypeByExtension(filepath.Ext(fileName))
		c.Writer.Header().Set("Content-Type", ctype)
		c.Writer.Write(data)
	})
}

func runGateWay(cmd *cobra.Command, args []string) error {
	cfgPath, err := utils.GetAbsPath(runConfigFile)
	if err != nil {
		return err
	}
	cfg, err := config.ParseConfig(cfgPath)
	if err != nil {
		return err
	}
	//steup gin
	engine := gin.New()
	engine.Use(gin.Recovery()) //TODO add authToken check
	apiRoute := engine.Group("/api")
	devRoute := apiRoute.Group("/dev")
	devapi.Init(devRoute, cfg)
	codingRoute := apiRoute.Group("/coding")
	codingapi.Init(codingRoute, cfg)
	initWeb(engine)
	//run server
	serverAddr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	if cfg.Ssl.Enable {
		err = engine.RunTLS(serverAddr, cfg.Ssl.Cert, cfg.Ssl.Key)
		if err != nil {
			return err
		}
	} else {
		err = engine.Run(serverAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	runCmd.Flags().StringVar(&runConfigFile, "config", "config.yaml", "set config file")
}

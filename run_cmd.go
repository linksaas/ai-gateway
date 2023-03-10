package main

import (
	"embed"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

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

func checkToken(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !strings.HasPrefix(ctx.Request.URL.Path, "/api") {
			ctx.Next()
			return
		}
		//api
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/dev/genToken") {
			ctx.Next()
			return
		}
		//check token
		token := ctx.GetHeader("X-AuthToken")
		contextValue, err := utils.DecodeToken(token, secret)
		if err != nil {
			utils.SendError(ctx, 403, err)
			ctx.Abort()
		} else {
			ctx.Request.Header.Set("contextValue", contextValue)
			ctx.Next()
		}
	}
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
	engine.Use(checkToken(cfg.Secret), gin.Recovery()) //TODO add authToken check
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

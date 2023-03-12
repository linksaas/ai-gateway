package main

import (
	"embed"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/arthurkiller/rollingwriter"
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

func initLogger(logDir string) (rollingwriter.RollingWriter, error) {
	if logDir == "" {
		logDir = "logs"
	}
	logDir, err := utils.GetAbsPath(logDir)
	if err != nil {
		return nil, err
	}
	os.MkdirAll(logDir, 0700) //skip error check
	return rollingwriter.NewWriterFromConfig(&rollingwriter.Config{
		LogPath:                logDir,
		TimeTagFormat:          "20060102",
		FileName:               "checker",
		MaxRemain:              30,
		RollingPolicy:          rollingwriter.TimeRolling,
		RollingTimePattern:     "0 0 0 * * *",
		RollingVolumeSize:      "100M",
		WriterMode:             "async",
		BufferWriterThershould: 8 * 1024 * 1024,
		Compress:               true,
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
	//init log
	logWriter, err := initLogger(cfg.LogDir)
	//steup gin
	engine := gin.New()
	engine.Use(checkToken(cfg.Secret), gin.Recovery()) //TODO add authToken check
	apiRoute := engine.Group("/api")
	devRoute := apiRoute.Group("/dev")
	devapi.Init(devRoute, cfg)
	codingRoute := apiRoute.Group("/coding")
	err = codingapi.Init(codingRoute, cfg, logWriter)
	if err != nil {
		return err
	}
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

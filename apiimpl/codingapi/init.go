package codingapi

import (
	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/config"
)

func Init(router *gin.RouterGroup, cfg *config.GateWayConfig) error {
	checker, err := newContentChecker(cfg)
	if err != nil {
		return err
	}
	scriptBackend, err := newScriptBackend(cfg)
	if err != nil {
		return err
	}

	router.POST("/complete/:lang", func(ctx *gin.Context) {
		handler := &CompleteHandler{}
		handler.process(ctx, cfg, checker, scriptBackend)
	})
	router.POST("/convert/:lang", func(ctx *gin.Context) {
		handler := &ConvertHandler{}
		handler.process(ctx, cfg, checker, scriptBackend)
	})
	router.POST("/explain/:lang", func(ctx *gin.Context) {
		handler := &ExplainHandler{}
		handler.process(ctx, cfg, checker, scriptBackend)
	})
	router.POST("/fixError/:lang", func(ctx *gin.Context) {
		handler := &FixErrorHandler{}
		handler.process(ctx, cfg, checker, scriptBackend)
	})
	router.POST("/genTest/:lang", func(ctx *gin.Context) {
		handler := &GenTestHandler{}
		handler.process(ctx, cfg, checker, scriptBackend)
	})
	return nil
}

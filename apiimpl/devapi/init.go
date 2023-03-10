package devapi

import (
	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/config"
)

func Init(router *gin.RouterGroup, cfg *config.GateWayConfig) {
	router.POST("/genToken", func(ctx *gin.Context) {
		handler := &GenTokenHandler{}
		handler.process(ctx, cfg.Secret, cfg.TokenTtl, cfg.Dev)
	})
	router.POST("/cap", func(ctx *gin.Context) {
		handler := &CapHandler{}
		handler.process(ctx, &cfg.Provider)
	})
}

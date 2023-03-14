package devapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/utils"
	ai_proto "github.com/linksaas/ai-proto-go"
)

type GenTokenRequest struct {
	ContextValue string `form:"contextValue"`
	RandomStr    string `form:"randomStr"`
}

type GenTokenHandler struct{}

func (handler *GenTokenHandler) process(ctx *gin.Context, secret string, tokenTtl int, dev bool) {
	if !dev {
		utils.SendError(ctx, 500, fmt.Errorf("not allowed in product mode"))
		return
	}
	req := &GenTokenRequest{}
	err := ctx.Bind(req)
	if err != nil {
		utils.SendError(ctx, 500, err)
		return
	}

	token, err := utils.GenToken(req.ContextValue, req.RandomStr, secret, tokenTtl)
	if err != nil {
		utils.SendError(ctx, 500, err)
		return
	}
	ctx.JSON(200, &ai_proto.ApiDevGenTokenPost200Response{
		Token: token,
	})
}

package devapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/utils"
	"github.com/linksaas/ai-proto-go/client"
)

type GenTokenHandler struct{}

func (handler *GenTokenHandler) process(ctx *gin.Context, secret string, tokenTtl int, dev bool) {
	if !dev {
		utils.SendError(ctx, 500, fmt.Errorf("not allowed in product mode"))
		return
	}
	reqBody := &client.ApiDevGenTokenPostRequest{}
	err := utils.ReadRequestBody(ctx, reqBody)
	if err != nil {
		utils.SendError(ctx, 500, err)
		return
	}
	token, err := utils.GenToken(reqBody.ContextValue, reqBody.RandomStr, secret, tokenTtl)
	if err != nil {
		utils.SendError(ctx, 500, err)
		return
	}
	ctx.JSON(200, &client.ApiDevGenTokenPost200Response{
		Token: token,
	})
}

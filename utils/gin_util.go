package utils

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-proto-go/client"
)

func ReadRequestBody(ctx *gin.Context, reqBody interface{}) error {
	defer ctx.Request.Body.Close()
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, reqBody)
	if err != nil {
		return err
	}
	return nil
}

func SendError(ctx *gin.Context, code int, err error) {
	errMsg := err.Error()
	ctx.JSON(code, &client.ErrInfo{
		ErrMsg: &errMsg,
	})
}

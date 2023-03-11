package utils

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	ai_proto "github.com/linksaas/ai-proto-go"
)

func SendError(ctx *gin.Context, code int, err error) {
	errMsg := err.Error()
	ctx.JSON(code, &ai_proto.ErrInfo{
		ErrMsg: &errMsg,
	})
}

func CopyResponse(ctx *gin.Context, response *http.Response) {
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		SendError(ctx, 500, err)
		return
	}
	ctx.Writer.WriteHeader(response.StatusCode)
	for k, v := range response.Header {
		ctx.Writer.Header()[k] = v
	}

	ctx.Writer.Write(data)
}

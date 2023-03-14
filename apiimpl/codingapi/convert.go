package codingapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-gateway/utils"
)

type ConvertHandler struct{}

func (handler *ConvertHandler) process(ctx *gin.Context, cfg *config.GateWayConfig, checker *ContentChecker, scriptBackend *ScriptBackend) {
	lang := ctx.Param("lang")
	destLang := ctx.Query("destLang")
	contentData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.SendError(ctx, 500, err)
		return
	}
	content := string(contentData)

	matchProviderIndex, done := checkContent(ctx, cfg.Provider.CodingProviderList, CODE_FUNC_CONVERT, lang, content, checker)
	if done {
		return
	}
	provider := cfg.Provider.CodingProviderList[matchProviderIndex]
	if strings.HasPrefix(provider.Backend, "script://") {
		script := scriptBackend.GetBackendScript(matchProviderIndex)
		if script == nil {
			utils.SendError(ctx, 500, fmt.Errorf("script not exist"))
			return
		}
		contentList, err := script.CallConvert(lang, destLang, content)
		if err != nil {
			utils.SendError(ctx, 500, err)
			return
		}
		ctx.JSON(200, contentList)
		return
	} else {
		backendUrl := fmt.Sprintf("%s%s?destLang=%s", provider.Backend, ctx.Request.URL.Path, destLang)
		res, err := http.Post(backendUrl, ctx.ContentType(), bytes.NewReader([]byte(content)))
		if err != nil {
			utils.SendError(ctx, 500, err)
			return
		}
		utils.CopyResponse(ctx, res)
	}
}

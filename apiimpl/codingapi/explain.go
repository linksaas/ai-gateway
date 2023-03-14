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

type ExplainHandler struct{}

func (handler *ExplainHandler) process(ctx *gin.Context, cfg *config.GateWayConfig, checker *ContentChecker, scriptBackend *ScriptBackend) {
	lang := ctx.Param("lang")
	contentData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.SendError(ctx, 500, err)
		return
	}
	content := string(contentData)

	matchProviderIndex, done := checkContent(ctx, cfg.Provider.CodingProviderList, CODE_FUNC_EXPLAIN, lang, content, checker)
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
		contentList, err := script.CallExplain(lang, content)
		if err != nil {
			utils.SendError(ctx, 500, err)
			return
		}
		ctx.JSON(200, contentList)
		return
	} else {
		backendUrl := fmt.Sprintf("%s%s", provider.Backend, ctx.Request.URL.Path)
		res, err := http.Post(backendUrl, ctx.ContentType(), bytes.NewReader([]byte(content)))
		if err != nil {
			utils.SendError(ctx, 500, err)
			return
		}
		utils.CopyResponse(ctx, res)
	}
}

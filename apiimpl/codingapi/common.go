package codingapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-gateway/utils"
)

type CODE_FUNC_TYPE int

const (
	CODE_FUNC_COMPLETE  = CODE_FUNC_TYPE(1)
	CODE_FUNC_CONVERT   = CODE_FUNC_TYPE(2)
	CODE_FUNC_EXPLAIN   = CODE_FUNC_TYPE(3)
	CODE_FUNC_FIX_ERROR = CODE_FUNC_TYPE(4)
	CODE_FUNC_GEN_TEST  = CODE_FUNC_TYPE(5)
)

func checkContent(ctx *gin.Context, providerList []config.CodingProviderConfig, funcType CODE_FUNC_TYPE, lang, content string, checker *ContentChecker) (matchProviderIndex int, done bool) {
	matchProviderIndex = -1
	for providerIndex, provider := range providerList {
		if len(provider.CompleteLangList) == 0 {
			continue
		}
		match := false
		cmpLangList := []string{}
		if funcType == CODE_FUNC_COMPLETE {
			cmpLangList = provider.CompleteLangList
		} else if funcType == CODE_FUNC_CONVERT {
			cmpLangList = provider.ConvertLangList
		} else if funcType == CODE_FUNC_EXPLAIN {
			cmpLangList = provider.ExplainLangList
		} else if funcType == CODE_FUNC_FIX_ERROR {
			cmpLangList = provider.FixerrorLangList
		} else if funcType == CODE_FUNC_GEN_TEST {
			cmpLangList = provider.GentestLangList
		}
		for _, cmpLang := range cmpLangList {
			if cmpLang == lang {
				match = true
				break
			}
		}
		if match {
			matchProviderIndex = providerIndex
		}
	}
	if matchProviderIndex == -1 {
		utils.SendError(ctx, 500, fmt.Errorf("not support lang %s", lang))
		done = true
		return
	}
	allow, err := checker.CheckCodeContent(ctx.Request.URL.Path, matchProviderIndex, content)
	if err != nil {
		utils.SendError(ctx, 403, err)
		done = true
		return
	}
	if !allow {
		utils.SendError(ctx, 403, fmt.Errorf("content not allow"))
		done = true
		return
	}
	return matchProviderIndex, false
}

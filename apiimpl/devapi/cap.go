package devapi

import (
	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/config"
	"github.com/linksaas/ai-proto-go/client"
)

type CapHandler struct{}

func (handler *CapHandler) process(ctx *gin.Context, cfg *config.ProviderConfig) {
	completeMap := map[string]bool{}
	genTestMap := map[string]bool{}
	convertMap := map[string]bool{}
	explainMap := map[string]bool{}
	fixErrorMap := map[string]bool{}
	for _, codingProvider := range cfg.CodingProviderList {
		for _, lang := range codingProvider.CompleteLangList {
			completeMap[lang] = true
		}
		for _, lang := range codingProvider.GentestLangList {
			genTestMap[lang] = true
		}
		for _, lang := range codingProvider.ConvertLangList {
			convertMap[lang] = true
		}
		for _, lang := range codingProvider.ExplainLangList {
			explainMap[lang] = true
		}
		for _, lang := range codingProvider.FixerrorLangList {
			fixErrorMap[lang] = true
		}
	}

	completeList := []client.Lang{}
	for k := range completeMap {
		completeList = append(completeList, client.Lang(k))
	}

	genTestList := []client.Lang{}
	for k := range genTestMap {
		genTestList = append(genTestList, client.Lang(k))
	}

	convertList := []client.Lang{}
	for k := range convertMap {
		convertList = append(convertList, client.Lang(k))
	}

	explainList := []client.Lang{}
	for k := range explainMap {
		explainList = append(explainList, client.Lang(k))
	}

	fixErrorList := []client.Lang{}
	for k := range fixErrorMap {
		fixErrorList = append(fixErrorList, client.Lang(k))
	}

	ctx.JSON(200, &client.ApiDevCapPost200Response{
		CompleteLangList: completeList,
		GenTestLangList:  genTestList,
		ConvertLangList:  convertList,
		ExplainLangList:  explainList,
		FixErrorLangList: fixErrorList,
	})
}

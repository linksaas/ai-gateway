package devapi

import (
	"github.com/gin-gonic/gin"
	"github.com/linksaas/ai-gateway/config"
	ai_proto "github.com/linksaas/ai-proto-go"
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

	completeList := []ai_proto.Lang{}
	for k := range completeMap {
		completeList = append(completeList, ai_proto.Lang(k))
	}

	genTestList := []ai_proto.Lang{}
	for k := range genTestMap {
		genTestList = append(genTestList, ai_proto.Lang(k))
	}

	convertList := []ai_proto.Lang{}
	for k := range convertMap {
		convertList = append(convertList, ai_proto.Lang(k))
	}

	explainList := []ai_proto.Lang{}
	for k := range explainMap {
		explainList = append(explainList, ai_proto.Lang(k))
	}

	fixErrorList := []ai_proto.Lang{}
	for k := range fixErrorMap {
		fixErrorList = append(fixErrorList, ai_proto.Lang(k))
	}

	ctx.JSON(200, &ai_proto.ApiDevCapPost200Response{
		Coding: ai_proto.ApiDevCapPost200ResponseCoding{
			CompleteLangList: completeList,
			GenTestLangList:  genTestList,
			ConvertLangList:  convertList,
			ExplainLangList:  explainList,
			FixErrorLangList: fixErrorList,
		},
	})
}

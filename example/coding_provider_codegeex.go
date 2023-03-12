package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CommonResponseResultOutput struct {
	Code []string `json:"code"`
}

type CommonResponseResult struct {
	Output CommonResponseResultOutput `json:"output"`
}

type CommonResponse struct {
	Message string               `json:"message"`
	Result  CommonResponseResult `json:"result"`
}

func adjustLang(lang string) string {
	if lang == "cplusplus" {
		return "C++"
	} else if lang == "csharp" {
		return "C#"
	}
	return strings.ToUpper(lang)[:1] + lang[1:]
}

// process /api/coding/complete/:lang
func Complete(lang, content string) []string {
	lang = adjustLang(lang)

	payload := map[string]interface{}{
		"n":      1,
		"lang":   lang,
		"prompt": content,
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req, err := http.NewRequest("POST", "https://wudao.aminer.cn/os/api/api/v2/multilingual_code/generate", bytes.NewBuffer(payloadData))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Origin", "https://codegeex.cn")
	req.Header.Add("Referer", "https://codegeex.cn/")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	response := &CommonResponse{}
	err = json.Unmarshal(resData, response)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return response.Result.Output.Code
}

// process /api/coding/convert/:lang
func Convert(lang, destLang, content string) []string {
	lang = adjustLang(lang)
	destLang = adjustLang(destLang)
	payload := map[string]interface{}{
		"src_lang": lang,
		"dst_lang": destLang,
		"prompt":   content,
		"n":        1,
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req, err := http.NewRequest("POST", "https://wudao.aminer.cn/os/api/api/v2/multilingual_code/translate", bytes.NewBuffer(payloadData))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Origin", "https://codegeex.cn")
	req.Header.Add("Referer", "https://codegeex.cn/")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(string(resData))
	response := &CommonResponse{}
	err = json.Unmarshal(resData, response)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return []string{strings.Join(response.Result.Output.Code, "")}
}

// process /api/coding/explain/:lang
func Explain(lang, content string) []string {
	lang = adjustLang(lang)
	payload := map[string]interface{}{
		"apikey":    "",
		"apisecret": "",
		"lang":      lang,
		"prompt":    content,
		"n":         1,
		"locale":    "zh-CN",
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req, err := http.NewRequest("POST", "https://wudao.aminer.cn/os/api/api/v2/multilingual_code/explain", bytes.NewBuffer(payloadData))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Origin", "https://codegeex.cn")
	req.Header.Add("Referer", "https://codegeex.cn/")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(string(resData))
	response := &CommonResponse{}
	err = json.Unmarshal(resData, response)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return []string{strings.Join(response.Result.Output.Code, "")}
}

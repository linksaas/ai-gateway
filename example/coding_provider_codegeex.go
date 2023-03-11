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

func runComplete(lang, content string) ([]string, error) {
	payload := map[string]interface{}{
		"n":      1,
		"lang":   lang,
		"prompt": content,
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://wudao.aminer.cn/os/api/api/v2/multilingual_code/generate", bytes.NewBuffer(payloadData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Origin", "https://codegeex.cn")
	req.Header.Add("Referer", "https://codegeex.cn/")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := &CommonResponse{}
	err = json.Unmarshal(resData, response)
	if err != nil {
		return nil, err
	}
	if response.Result.Output.Code == nil {
		return nil, fmt.Errorf("error response")
	}
	return response.Result.Output.Code, nil
}

// process /api/coding/complete/:lang
func Complete(lang, content string) []string {
	lang = strings.ToUpper(lang)[:1] + lang[1:]
	lines := strings.Split(content, "\n")
	firstLine := ""
	for _, line := range lines {
		line = strings.Trim(line, " \r\n\t")
		if line != "" {
			firstLine = line
			break
		}
	}
	if firstLine == "" {
		return nil
	}

	allResult := []string{}
	stepCount := 0
	for {
		stepCount++
		if stepCount > 100 {
			break
		}
		lines, err := runComplete(lang, content)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if len(lines) == 0 {
			break
		}
		// fmt.Println(lines[0])
		content += "\n"
		content += lines[0]
		allResult = append(allResult, lines[0])
		if strings.Contains(lines[0], firstLine) {
			break
		}
		if strings.Trim(lines[0], " \r\n\t") == "" {
			break
		}
	}
	return []string{strings.Join(allResult, "")}
}

// process /api/coding/convert/:lang
func Convert(lang, destLang, content string) []string {
	lang = strings.ToUpper(lang)[:1] + lang[1:]
	destLang = strings.ToUpper(destLang)[:1] + destLang[1:]
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
	lang = strings.ToUpper(lang)[:1] + lang[1:]
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

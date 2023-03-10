package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TokenInfo struct {
	ContextValue string `yaml:"contextValue"`
	RandomStr    string `yaml:"randomStr"`
	ExpireTime   int64  `yaml:"expireTime"`
	SignStr      string `yaml:"signStr"`
}

func GenToken(contextValue, randomStr, secret string, ttl int) (string, error) {
	if len(randomStr) < 32 {
		return "", fmt.Errorf("random string must have 32 chars at least")
	}
	expireTime := int64(-1)
	if ttl > 0 {
		expireTime = time.Now().Unix() + int64(ttl)
	}
	toSignStr := fmt.Sprintf("%s%d%s%s", contextValue, expireTime, randomStr, secret)
	h := sha256.New()
	h.Write([]byte(toSignStr))
	signStr := hex.EncodeToString(h.Sum(nil))
	signStr = strings.ToUpper(signStr)

	tokenInfo := &TokenInfo{
		ContextValue: contextValue,
		RandomStr:    randomStr,
		ExpireTime:   expireTime,
		SignStr:      signStr,
	}
	tokenData, err := json.Marshal(tokenInfo)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenData), nil
}

func DecodeToken(tokenStr, secret string) (string, error) {
	tokenData, err := hex.DecodeString(tokenStr)
	if err != nil {
		return "", err
	}
	tokenInfo := &TokenInfo{}
	err = json.Unmarshal(tokenData, tokenInfo)
	if err != nil {
		return "", err
	}
	if tokenInfo.ExpireTime > 0 && time.Now().Unix() > tokenInfo.ExpireTime {
		return "", fmt.Errorf("token expire")
	}
	toSignStr := fmt.Sprintf("%s%d%s%s", tokenInfo.ContextValue, tokenInfo.ExpireTime, tokenInfo.RandomStr, secret)
	h := sha256.New()
	h.Write([]byte(toSignStr))
	signStr := hex.EncodeToString(h.Sum(nil))
	signStr = strings.ToUpper(signStr)
	if signStr == tokenInfo.SignStr {
		return tokenInfo.ContextValue, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

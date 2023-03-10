package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func GetAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	curPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	curPath = filepath.Dir(curPath)
	return strings.Join([]string{curPath, path}, string(os.PathSeparator)), nil
}

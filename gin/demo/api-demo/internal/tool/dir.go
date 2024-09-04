package tool

import (
	"os"
	"path/filepath"
)

// 获取项目根目录
func GetRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return rootDir, nil
}

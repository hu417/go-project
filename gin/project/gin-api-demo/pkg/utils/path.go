package utils

import (
	"os"
	"path"
)

// 判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 获取当前路径
func GetPath() string {
	dir, _ := os.Getwd()
	// 使用path包的Join方法获取当前目录的路径
	return path.Join(dir)
}

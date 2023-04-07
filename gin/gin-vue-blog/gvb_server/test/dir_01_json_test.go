package test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 1、根据json文件创建目录结构 //

var (
	stRootDir   string // 定义根目录目录
	stSeparator string // 定义路径分隔符
	stJsonFile  string = "dir.json"
	iJsonData   map[string]any
)

// 获取目录/文件信息
func loadJson() {
	// 当前目录
	stWorkDir, _ := os.Getwd()
	// 路径分隔符
	stSeparator = string(filepath.Separator)
	// 设置根目录
	stRootDir = stWorkDir[:strings.LastIndex(stWorkDir, stSeparator)]
	fmt.Println(stRootDir) // 获取根目录

	// 获取json文件全路径
	gnJsonBytes, _ := os.ReadFile(stWorkDir + stSeparator + stJsonFile)
	// 读取json文件内容
	err := json.Unmarshal(gnJsonBytes, &iJsonData)
	if err != nil {
		panic(err)
	}
}

// 解析json
func parseJsonMap(mapData map[string]any, stParetentDir string) {
	for k, v := range mapData {

		// 判断 v的数据类型
		switch v := v.(type) {
		case string:
			{
				path := v // 类型转换
				if path == "" {
					continue
				}

				if stParetentDir != "" {
					path = stParetentDir + stSeparator + path
					if k == "text" {
						stParetentDir = path
					}

				} else {
					stParetentDir = path
				}
				createDir(path)
			}
		case []any:
			{
				parseArray(v, stParetentDir)
			}

		}
	}
}

func parseArray(giJsonData []any, stParentDir string) {
	for _, v := range giJsonData {
		mapV, _ := v.(map[string]any)
		parseJsonMap(mapV, stParentDir)
	}
}

// 创建目录
func createDir(path string) {
	if path == "" {
		return
	}
	fmt.Println(path)
	// 创建目录
	err := os.MkdirAll(stRootDir+stSeparator+path, fs.ModePerm)
	if err != nil {
		panic(err)
	}

}

func TestGenerlDir(t *testing.T) {
	loadJson()
	// 解析json
	parseJsonMap(iJsonData, "")
}

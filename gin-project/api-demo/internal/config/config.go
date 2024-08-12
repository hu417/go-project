package config

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

// Config 配置文件集合
type Config struct {
	App      App
	Database Database
	Logger   Logger
}

// GetConfig 读取配置文件
func GetConfig() (c *Config) {
	// 获取项目的根目录
	// rootDir, _ := tool.GetRootDir()

	// 获取项目的根目录
	_, filePath, _, _ := runtime.Caller(0)
	dirPath := path.Dir(filePath)
	rootDir := dirPath + "/../../"

	// 实例化viper，并根据地址读取配置文件
	v := viper.New()
	v.SetConfigFile(rootDir + "/etc/conf.yaml")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("ReadInConfigError: ", err)
		return
	}

	// 将读取到的配置文件绑定到返回参数c
	err = v.Unmarshal(&c)
	if err != nil {
		fmt.Println("ConfigUnmarshalError: ", err)
		return
	}

	return
}

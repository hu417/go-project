package core

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"gvb_server/config"
	"gvb_server/global"

	"gopkg.in/yaml.v2"
)

const ConfigFile = "./service/settings.yaml"

// 读取yaml文件的配置信息
func InitConf() {

	c := &config.Config{}
	yamlConf, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error: %v", err))

	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		panic(err)

	}

	global.Config = c
}

// 写入yaml文件的配置信息
func WriteConf() error {

	// 获取配置信息
	byteDate, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}

	// 写入配置信息
	err = ioutil.WriteFile(ConfigFile, byteDate, fs.ModePerm)
	if err != nil {
		return err
	}
	global.Logger.Info("配置修改成功")
	return nil
}

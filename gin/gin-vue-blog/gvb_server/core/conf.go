package core

import (
	"fmt"
	"io/ioutil"

	"gvb_server/config"
	"gvb_server/global"

	"gopkg.in/yaml.v2"
)

// 读取yaml文件的配置信息
func InitConf() {

	const ConfigFile = "./service/settings.yaml"
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

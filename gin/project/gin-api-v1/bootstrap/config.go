package bootstrap

import (
	"fmt"

	"gin-api-demo/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"os"
)

func InitConfig(config string) *viper.Viper {
	// 设置配置文件路径
	if config == "" {
		config = "config.yaml"
	}

	// 生产环境可以通过设置环境变量来改变配置文件路径
	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		config = configEnv
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed:: %s", err))
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 重载配置
		if err := v.Unmarshal(&global.Conf); err != nil {
			panic(fmt.Errorf("config unmarshal failed:: %s", err))
		}
	})
	// 将配置赋值给全局变量
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Errorf("config unmarshal failed:: %s", err))
	}

	return v
}

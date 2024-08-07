package initialize

import (
	"flag"
	"fmt"
	"os"

	"gin-base/config"
	"gin-base/global"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitializeViper 优先级: 命令行 > 环境变量 > 默认值
func InitViper(path ...string) *config.Config {
	var conf_file string

	if len(path) == 0 {
		// 定义命令行flag参数，格式：flag.TypeVar(Type指针, flag名, 默认值, 帮助信息)
		flag.StringVar(&conf_file, "c", "", "choose config file.")

		// 定义好命令行flag参数后，需要通过调用flag.Parse()来对命令行参数进行解析。
		flag.Parse()

		// 判断命令行参数是否为空
		if conf_file == "" {
			/*
			   判断 global.ConfigEnv 常量存储的环境变量是否为空
			   比如我们启动项目的时候，执行：GVA_CONFIG=config.yaml go run main.go
			   这时候 os.Getenv(global.ConfigEnv) 得到的就是 config.yaml
			   当然，也可以通过 os.Setenv(global.ConfigEnv, "config.yaml") 在初始化之前设置
			*/
			if configEnv := os.Getenv(global.ConfigEnv); configEnv == "" {
				switch gin.Mode() {
				case gin.DebugMode:
					conf_file = global.ConfigDefaultFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, global.ConfigDefaultFile)
				case gin.ReleaseMode:
					conf_file = global.ConfigReleaseFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, global.ConfigReleaseFile)
				case gin.TestMode:
					conf_file = global.ConfigTestFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, global.ConfigTestFile)
				}
			} else {
				// global.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				conf_file = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", global.ConfigEnv, conf_file)
			}
		} else {
			// 命令行参数不为空 将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", conf_file)
		}
	} else {
		// 函数传递的可变参数的第一个值赋值于config
		conf_file = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", conf_file)
	}

	vip := viper.New()
	vip.SetConfigFile(conf_file)
	vip.SetConfigType("yaml")

	if err := vip.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	conf := &config.Config{}
	vip.WatchConfig()

	vip.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := vip.Unmarshal(conf); err != nil {
			fmt.Println(err)
		}
	})

	if err := vip.Unmarshal(conf); err != nil {
		fmt.Println(err)
	}

	fmt.Println("====1-viper====: viper init config success")

	return conf
}

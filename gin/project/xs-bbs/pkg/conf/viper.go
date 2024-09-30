package conf

// 初始化viper
import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitViper(file string) *Config {
	// 指定配置文件
	// 方式1：直接指定配置文件路径（相对路径或者绝对路径）
	// 相对路径：相对执行的可执行文件的相对路径
	viper.SetConfigFile(file)
	// 指定在什么路径下查找配置文件(相对路径)
	//viper.AddConfigPath("./conf")
	// 读取配置文件信息
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 反序列化到Conf中
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	// 持续监视配置文件是否发生变化
	viper.WatchConfig()
	// 配置文件发生改变执行的回调函数
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了")
		// 重新反序列化到Conf中
		if err := viper.Unmarshal(&conf); err != nil {
			panic(fmt.Sprintf("viper.Unmarshal failed, err:%v\n", err))
		}
	})
	return &conf
}

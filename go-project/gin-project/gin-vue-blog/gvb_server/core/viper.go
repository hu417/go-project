package core

// 安装依赖: go get github.com/spf13/viper
import (
	"github.com/spf13/viper"
)

func IninViper() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//fmt.Println(viper.GetString("Server.Name"))
}

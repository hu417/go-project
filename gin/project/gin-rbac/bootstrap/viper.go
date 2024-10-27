package bootstrap

import (
	"fmt"
	"os"

	"gin-rbac/config"

	"github.com/spf13/viper"
)

// LoadConfig 从配置文件中加载配置
func LoadConfig() (*config.Config, error) {
	// 设置配置文件名
	viper.SetConfigName("setting")
	// 设置配置文件路径
	viper.AddConfigPath("./gin-rbac/etc")
	// viper.AddConfigPath("./etc")

	// 读取环境变量，如果存在，则使用环境特定的配置文件
	env := os.Getenv("ENV")
	if env != "" {
		viper.SetConfigFile(fmt.Sprintf("settings/settings.%s.yaml", env))
	}

	// 查找并读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// 将配置文件的内容解析（映射）到Config结构体中
	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

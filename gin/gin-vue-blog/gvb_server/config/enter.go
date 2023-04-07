package config

// 此文件表示是要做的第一件事情
type Config struct {
	Mysql  Mysql  `yaml:"mysql"`
	System System `yaml:"system"`
	Logger Logger `yaml:"logger"`
}

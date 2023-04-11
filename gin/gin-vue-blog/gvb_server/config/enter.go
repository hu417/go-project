package config

// 此文件表示是要做的第一件事情
type Config struct {
	Mysql    Mysql    `yaml:"mysql"`
	System   System   `yaml:"system"`
	Logger   Logger   `yaml:"logger"`
	Elastic  Elastic  `yaml:"elastic"`
	SiteInfo SiteInfo `yaml:"site_info"`
	Email    Email    `yaml:"email"`
	QQ       QQ       `yaml:"qq"`
	QiNiu    QiNiu    `yaml:"qi_niu"`
	JWT      JWT      `yaml:"jwt"`
}

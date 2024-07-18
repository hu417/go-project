package config

type Upload struct {
	// 文件上传
	Path string `json:"path" yaml:"path"`
	Size int    `json:"size" yaml:"size"`
}

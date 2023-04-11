package config

// 七牛云 对象存储配置
type QiNiu struct {
	AccessKey string  `json:"access_key" yaml:"access_key"` // api token
	SecretKey string  `json:"secret_key" yaml:"secret_key"` // api secret
	Bucket    string  `json:"bucket" yaml:"bucket"`         // 存储桶的名字
	CDN       string  `json:"cdn" yaml:"cdn"`               // 访问图片的地址的前缀
	Zone      string  `json:"zone" yaml:"zone"`             // 存储的地区
	Size      float64 `json:"size" yaml:"size"`             // 存储的大小限制，单位是MB
}

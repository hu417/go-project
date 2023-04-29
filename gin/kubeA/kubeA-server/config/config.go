package config

const (
	// 定义服务监听地址
	ListenAddr = "0.0.0.0:9090"

	// 添加kubeconfig文件
	// .表示相对于项目的根目录,也可以写绝对路径
	Kubeconfigs = `{
					"Test1": "./file/config1",
					"Test2": "./file/config2"
					}`
)

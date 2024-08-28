package v1

type Registration struct {
	Name    string   `json:"name"`    // 服务名称
	ID      string   `json:"id"`      // 服务 ID，必须唯一
	Address string   `json:"address"` // 服务的地址
	Port    int      `json:"port"`    // 服务端口 服务所在的监听端口
	Tags    []string `json:"tags"`    // 可选：服务标签
}

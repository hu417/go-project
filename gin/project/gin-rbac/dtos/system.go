package dtos

type RouterDTO struct {
	Method  string `json:"method"`   // 请求方法
	ApiPath string `json:"api_path"` // 接口路径
}

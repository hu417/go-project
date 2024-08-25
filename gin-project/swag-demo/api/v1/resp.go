package v1

import "swag-demo/model"

// 响应结构体
type RespuserList struct {
	Code    int           `json:"code"`    // 业务响应状态码
	Message string        `json:"message"` // 提示信息
	Data    []*model.User `json:"data"`    // 数据
}

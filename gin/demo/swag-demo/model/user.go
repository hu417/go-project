package model

type User struct {
	// 设置最小值、最大值、默认值
	ID int `json:"id" minimum:"10" maximum:"20" default:"15"`
	// 设置最小长度、最大长度、示例值
	Name string `json:"name" maxLength:"16" example:"random string"`
	Age  int    `json:"age"`
}

package api

// User 用户结构体
type User struct {
	Name    string `json:"name" bson:"name,omitempty"`
	Age     int    `json:"age" bson:"age,omitempty"`
	Address string `json:"address" bson:"address,omitempty"`
}

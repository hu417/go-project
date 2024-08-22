package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应string数据
func RespStr(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

// 响应json数据
func RespJson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": nil,
	})
}

// 响应yaml数据
func RespYaml(c *gin.Context) {
	c.YAML(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": nil,
	})
}

// 响应xml数据
func RespXml(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": nil,
	})
}

// 响应protobuf数据
func RespProtobuf(c *gin.Context) {
	// reps := []int64{1, 2}
	// data := &protoexample.Test{
	// 	Reps: reps,
	// }
	c.ProtoBuf(200, "ok")
}

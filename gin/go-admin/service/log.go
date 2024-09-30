package service

import (
	"go-admin/api/request"
	"go-admin/dao/db"

	"github.com/gin-gonic/gin"
)

// GetLogList 获取日志列表
func GetLogList(c *gin.Context, in *request.GetLogListRequest) (interface{}, error) {

	var (
		cnt  int64
		list = make([]*request.GetLogListReply, 0)
	)
	err := db.GetLogList(in.Keyword).Count(&cnt).Offset((in.Page - 1) * in.Size).Limit(in.Size).Find(&list).Error

	data := struct {
		List  []*request.GetLogListReply `json:"list"`
		Count int64                      `json:"count"`
	}{
		List:  list,
		Count: cnt,
	}

	return data, err
}

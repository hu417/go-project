package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getPageInfo(c *gin.Context) (int64, int64) {
	// 分页参数的处理
	var (
		page int64
		size int64
		err  error
	)

	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}

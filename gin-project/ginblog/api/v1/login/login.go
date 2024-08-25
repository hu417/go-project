package login

import (
	v1 "ginblog/api/v1"
	"ginblog/global"
	"ginblog/service/user"
	"ginblog/utils/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 管理员登录
func LoginByAdmin(ctx *gin.Context) {

}

// 前端登录
func LoginByFront(ctx *gin.Context) {
	var u v1.UserParams
	// 绑定参数
	err := ctx.ShouldBind(&u)
	if err != nil {
		resp.Fail(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 检测用户名是否已存在
	ok := user.NewUserSvc(global.DB).CheckUser(u.Username)

	if !ok {
		resp.Fail(ctx, http.StatusBadRequest, "用户名不存在", nil)
		return
	}

	resp.Success(ctx, 200, "用户名已存在", nil)
}

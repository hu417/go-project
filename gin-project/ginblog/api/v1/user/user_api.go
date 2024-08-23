package user

import "github.com/gin-gonic/gin"

type User struct {
}

type UserApiInterface interface {
	Add(ctx *gin.Context)
}

func NewUser() *User {
	return &User{}
}

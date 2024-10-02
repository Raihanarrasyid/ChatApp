package user

import "github.com/gin-gonic/gin"

type UserController interface {
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
}
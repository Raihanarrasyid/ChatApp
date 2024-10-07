package auth

import "github.com/gin-gonic/gin"

type AuthController interface {
	RequestOtpSignUp(ctx *gin.Context)
	VerifyOtpSignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
}
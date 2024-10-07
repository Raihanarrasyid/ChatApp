package auth

import (
	"ChatApp/configs"
	"ChatApp/internal/http/request"
	"ChatApp/internal/service/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthControllerImpl struct {
	authService auth.AuthService
	config      *configs.Config
}

func NewAuthController(router *gin.RouterGroup, authService auth.AuthService, config *configs.Config) {
	controller := &AuthControllerImpl{
		authService: authService,
		config:      config,
	}

	router.POST("/signup/request-otp", controller.RequestOtpSignUp)
	router.POST("/signup/verify-otp", controller.VerifyOtpSignUp)
	router.POST("/signin", controller.SignIn)
}

//	@Summary		Request OTP for sign up
//	@Description	Request OTP for sign up
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.CreateUserRequest	true	"Create User"
//	@Success		200		{object}	http.Response
//	@Failure		400		{object}	http.Error
//	@Failure		500		{object}	http.Error
//	@Router			/auth/signup/request-otp [post]
func (auc *AuthControllerImpl) RequestOtpSignUp(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}	
	err := auc.authService.SendOTP(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP sent",
	})
}

//	@Summary		Verify OTP for sign up
//	@Description	Verify OTP for sign up
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.VerifyOtpRequest	true	"Verify OTP"
//	@Success		200		{object}	http.Response
//	@Failure		400		{object}	http.Error
//	@Failure		500		{object}	http.Error
//	@Router			/auth/signup/verify-otp [post]
func (auc *AuthControllerImpl) VerifyOtpSignUp(ctx *gin.Context) {
	var req request.VerifyOtpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := auc.authService.VerifyAndCreateUser(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP verified",
	})
}

//	@Summary		Sign in
//	@Description	Sign in
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.SignInRequest	true	"Sign In"
//	@Success		200		{object}	http.Response
//	@Failure		400		{object}	http.Error
//	@Failure		401		{object}	http.Error
//	@Failure		500		{object}	http.Error
//	@Router			/auth/signin [post]
func (auc *AuthControllerImpl) SignIn(ctx *gin.Context) {
	var req request.SignInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, accessToken, refreshToken, err := auc.authService.SignIn(ctx, req, auc.config.JwtSecret)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
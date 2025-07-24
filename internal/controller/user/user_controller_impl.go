package user

import (
	"net/http"
	"strconv"

	"ChatApp/internal/http/request"
	"ChatApp/internal/service/user"

	"github.com/gin-gonic/gin"
)

type UserControllerImpl struct {
	UserService user.UserService
}

func NewUserController(router *gin.RouterGroup, userService user.UserService) {
	controller := &UserControllerImpl{
		UserService: userService,
	}

	router.GET("/", controller.GetAllUser)
	router.PUT("/", controller.UpdateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.GET("/:id", controller.GetUserByID)
}

//	@Summary		Update user
//	@Description	Update user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID"
//	@Param			body	body		request.UpdateUserRequest	true	"Update User"
//	@Success		200		{object}	http.Response
//	@Failure		400		{object}	http.Error
//	@Failure		404		{object}	http.Error
//	@Failure		500		{object}	http.Error
//	@Router			/users [put]
func (u *UserControllerImpl) UpdateUser(ctx *gin.Context) {
	id := ctx.GetString("user_id")
	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := u.UserService.Update(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated",
	})
}

//	@Summary		Delete user
//	@Description	Delete user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	http.Response
//	@Failure		400	{object}	http.Error
//	@Failure		404	{object}	http.Error
//	@Failure		500	{object}	http.Error
//	@Router			/users [delete]
func (u *UserControllerImpl) DeleteUser(ctx *gin.Context) {
	id := ctx.GetString("user_id")
	err := u.UserService.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
	})
}

//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	http.Response{value=response.UserResponse}
//	@Failure		400	{object}	http.Error
//	@Failure		404	{object}	http.Error
//	@Failure		500	{object}	http.Error
//	@Router			/users/{id} [get]
func (u *UserControllerImpl) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := u.UserService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

//	@Summary		Get all user
//	@Description	Get all user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int	false	"Page"
//	@Param			page_size	query		int	false	"Page Size"
//	@Success		200			{object}	http.Response{value=[]response.UserResponse}
//	@Failure		400			{object}	http.Error
//	@Failure		500			{object}	http.Error
//	@Router			/users [get]
func (u *UserControllerImpl) GetAllUser(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page",
		})
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page size",
		})
		return
	}

	users, err := u.UserService.GetAll(ctx, page, pageSize)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
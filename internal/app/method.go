package app

import (
	http "ChatApp/internal/http/server"
	pg "ChatApp/pkg/db"
	"log"

	UserRepository "ChatApp/internal/repository/user"

	UserService "ChatApp/internal/service/user"

	UserController "ChatApp/internal/controller/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (app *App) Run() {
	server := http.NewHttpServer(app.config.GinMode)

	server.GET("", func (ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	var err error
	var db *gorm.DB
	db, err = pg.NewDB(app.config.DBHost)

	if err != nil {
		panic(err)
	}

	pg.Migrate(db)

	initControllers(server, db)

	log.Printf("Server is running on port %s mode %s", app.config.Port, app.config.GinMode)
	err = server.Run(":" + app.config.Port)

	if err != nil {
		panic(err)
	}
}

func initControllers(
	router *gin.Engine,
	db *gorm.DB,
) {
	userRepository := UserRepository.NewUserRepository(db)

	userService := UserService.NewUserService(userRepository)

	routerGroup := router.Group("/api/v1")


	UserController.NewUserController(routerGroup.Group("/users"), userService)
}
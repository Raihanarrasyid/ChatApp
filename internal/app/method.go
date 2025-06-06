package app

import (
	"ChatApp/configs"
	"ChatApp/docs"
	http "ChatApp/internal/http/server"
	"ChatApp/internal/middleware"
	database "ChatApp/pkg/db"
	"log"

	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	AuthRepository "ChatApp/internal/repository/auth"
	ChatRepository "ChatApp/internal/repository/chat"
	UserRepository "ChatApp/internal/repository/user"

	AuthService "ChatApp/internal/service/auth"
	ChatService "ChatApp/internal/service/chat"
	EmailService "ChatApp/internal/service/email"
	UserService "ChatApp/internal/service/user"

	AuthController "ChatApp/internal/controller/auth"
	ChatController "ChatApp/internal/controller/chat"
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
	var postgresDB *gorm.DB
	var redisDB *redis.Client
	postgresDB, err = database.NewPostgresDB(app.config.DBHost)
	redisDB, err = database.NewRedisDB(app.config.RedisHost, app.config.RedisPass)

	if err != nil {
		panic(err)
	}

	database.Migrate(postgresDB)

	if app.config.GinMode == "debug" {
		log.Println("Init Swagger")
		docs.SwaggerInfo.Version = "1.0.0"
		docs.SwaggerInfo.Title = "ChatApp API"
		docs.SwaggerInfo.Description = "ChatApp API"
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}


		server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:3000/swagger/doc.json")))
	}

	initControllers(server, postgresDB, redisDB, app.config)

	log.Printf("Server is running on port %s mode %s", app.config.Port, app.config.GinMode)
	err = server.Run(":" + app.config.Port)

	if err != nil {
		panic(err)
	}
}

func initControllers(
	router *gin.Engine,
	postgresDB *gorm.DB,
	redisDB *redis.Client,
	config *configs.Config,
) {
	userRepository := UserRepository.NewUserRepository(postgresDB)
	authRepository := AuthRepository.NewAuthRepositoryImpl(redisDB)
	chatRepository := ChatRepository.NewChatRepository(redisDB, postgresDB)

	userService := UserService.NewUserService(userRepository)
	emailService := EmailService.NewEmailService(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	authService := AuthService.NewAuthService(userRepository, authRepository, *emailService)
	chatService := ChatService.NewChatService(chatRepository)

	
	routerGroup := router.Group("/api/v1")

	protected := routerGroup.Group("/")
	protected.Use(middleware.AuthMiddleware(config))
	
	AuthController.NewAuthController(routerGroup.Group("/auth"), authService, config)
	UserController.NewUserController(protected.Group("/users"), userService)
	ChatController.NewChatController(protected.Group("/chat"), chatService)
}
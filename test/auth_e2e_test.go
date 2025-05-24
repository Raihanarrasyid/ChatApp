package test

import (
	"ChatApp/configs"
	AuthController "ChatApp/internal/controller/auth"
	"ChatApp/internal/http/request"
	AuthRepository "ChatApp/internal/repository/auth"
	UserRepository "ChatApp/internal/repository/user"
	AuthService "ChatApp/internal/service/auth"
	EmailService "ChatApp/internal/service/email"
	database "ChatApp/pkg/db"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		DB:   1,               
	})

	_ = rdb.FlushDB(context.Background())

	config := &configs.Config{
		JwtSecret: "testsecret",
	}

	emailSvc := &EmailService.EmailService{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "465",
		Sender:   "",
		Password: "",
	}

	emailService := EmailService.NewEmailService(emailSvc.SMTPHost, emailSvc.SMTPPort, emailSvc.Sender, emailSvc.Password)

	gormDB, _ := database.NewPostgresDB("postgres://postgres:1234test@localhost:5432/chat_app")

	userRepository := UserRepository.NewUserRepository(gormDB)
	authRepo := AuthRepository.NewAuthRepositoryImpl(rdb)
	authSvc := AuthService.NewAuthService(userRepository, authRepo, *emailService)
	router := gin.Default()

	v1 := router.Group("/auth")
	AuthController.NewAuthController(v1, authSvc, config)

	return router
}

func TestRequestOTP(t *testing.T) {
	router := setupTestServer()

	reqBody := request.CreateUserRequest{
		Email:    "",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup/request-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestVerifyOTP(t *testing.T) {
	router := setupTestServer()

	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
	redisClient.Set(ctx, "user@example.com", "123456", 5*time.Minute)

	verifyReq := request.VerifyOtpRequest{
		Email: "user@example.com",
		OTP:   "123456",
	}

	body, _ := json.Marshal(verifyReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup/verify-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}


func TestSignIn(t *testing.T) {
	router := setupTestServer()

	signinReq := request.SignInRequest{
		Email:    "user@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(signinReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusUnauthorized)
}
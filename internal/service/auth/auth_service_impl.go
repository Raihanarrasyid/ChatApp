package auth

import (
	"ChatApp/internal/http/request"
	"ChatApp/internal/http/response"
	"ChatApp/internal/model"
	"ChatApp/internal/repository/auth"
	"ChatApp/internal/repository/user"
	"ChatApp/internal/service/email"
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	commonError "ChatApp/internal/util/error"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	UserRepository user.UserRepository
	AuthRepository auth.AuthRepository
	EmailService email.EmailService
}

func NewAuthService(userRepository user.UserRepository, authRepository auth.AuthRepository, emailService email.EmailService) AuthService {
	return &AuthServiceImpl{userRepository, authRepository, emailService}
}

func (a *AuthServiceImpl) SendOTP(ctx context.Context, req request.CreateUserRequest) error {
	otp := generateOTP()
	err := a.AuthRepository.SaveOTP(ctx, req.Email, otp, 5*time.Minute)
	if err != nil {
		return err
	}

	err = a.EmailService.SendOTP(req.Email, otp)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceImpl) VerifyAndCreateUser(ctx context.Context, req request.VerifyOtpRequest) error {
	isValid := a.AuthRepository.VerifyOTP(ctx, req.Email, req.OTP)
	if !isValid {
		return fmt.Errorf("invalid otp")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return commonError.NewBadRequest("Failed to hash password: " + err.Error())
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := a.UserRepository.Create(ctx, &user); err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceImpl) SignIn(ctx context.Context, req request.SignInRequest, jwtSecret string) (signin response.SignInResponse, err error) {
	user, err := a.UserRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return signin, commonError.NewNotFound("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return signin, commonError.NewBadRequest("Invalid password")
	}

	accessToken, err := generateAccessToken(user, jwtSecret)
	if err != nil {
		return signin, err
	}

	refreshToken, err := generateRefreshToken(user, jwtSecret)
	if err != nil {
		return signin, err
	}

	err = a.AuthRepository.SaveRefreshToken(ctx, user.ID.String(), refreshToken, 30*24*time.Hour)
	if err != nil {
		return signin, err
	}

	userResponse := response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	signin = response.SignInResponse{
		User:          userResponse,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
	}

	return signin, nil
}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.IntN(1000000))
}

func generateAccessToken(user *model.User, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id" : user.ID,
		"exp" : time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func generateRefreshToken(user *model.User, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id" : user.ID,
		"exp" : time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}


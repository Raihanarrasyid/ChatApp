package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthRepositoryImpl struct {
	redisDB *redis.Client
}

func NewAuthRepositoryImpl(redisDB *redis.Client) AuthRepository {
	return &AuthRepositoryImpl{redisDB}
}

func (a *AuthRepositoryImpl) SaveOTP(ctx context.Context, email, otp string, expiryTime time.Duration) error {
	err := a.redisDB.Set(ctx, email, otp, expiryTime).Err()
	if err != nil {
		return err
	}
	return nil	
}

func (a *AuthRepositoryImpl) VerifyOTP(ctx context.Context, email, otp string) bool {
	storedOTP, err := a.redisDB.Get(ctx, email).Result()
	if err != nil {
		return false
	}

	if storedOTP == otp {
		a.redisDB.Del(ctx, email)
		return true
	}

	return false
}

func (rr *AuthRepositoryImpl) SaveRefreshToken(ctx context.Context, userID string, refreshToken string, expiry time.Duration) error {
	err := rr.redisDB.Set(ctx, "refresh_token_"+userID, refreshToken, expiry).Err()
	if err != nil {
		return fmt.Errorf("failed to save refresh token in Redis: %v", err)
	}
	return nil
}

func (rr *AuthRepositoryImpl) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := rr.redisDB.Get(ctx, "refresh_token_"+userID).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}
package auth

import (
	"context"
	"time"
)

type AuthRepository interface {
	SaveOTP(ctx context.Context, email, otp string, expiryTime time.Duration) error
	VerifyOTP(ctx context.Context, email, otp string) bool
	SaveRefreshToken(ctx context.Context, email, refreshToken string, expiryTime time.Duration) error
	GetRefreshToken(ctx context.Context, email string) (string, error)
}
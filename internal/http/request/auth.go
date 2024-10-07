package request

type CreateUserRequest struct {
	Email string `json:"email"`
}

type VerifyOtpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
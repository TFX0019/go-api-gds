package auth

import (
	"errors"
	"time"

	"github.com/TFX0019/api-go-gds/pkg/config"
	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, string, *UserResponse, error)
	VerifyEmail(token string) error
	RefreshToken(tokenString string) (string, error)
	ForgotPassword(req ForgotPasswordRequest) error
	VerifyCode(req VerifyCodeRequest) error
	ResetPassword(req ResetPasswordRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(req RegisterRequest) error {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	verificationToken := uuid.New().String()

	user := &User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hashedPassword,
		VerificationToken: verificationToken,
		IsVerified:        false,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	utils.SendVerificationEmail(user.Email, verificationToken)
	return nil
}

func (s *service) Login(req LoginRequest) (string, string, *UserResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", "", nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", "", nil, errors.New("please verify your email first")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	userResponse := &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return accessToken, refreshToken, userResponse, nil
}

func (s *service) VerifyEmail(token string) error {
	user, err := s.repo.FindVerifyToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	return s.repo.UpdateUser(user)
}

func (s *service) RefreshToken(tokenString string) (string, error) {
	secret := config.GetEnv("JWT_REFRESH_SECRET", "refresh_secret")
	token, err := utils.ValidateToken(tokenString, secret)
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("invalid user id in token")
	}
	userID := uint(userIDFloat)

	// In a strict implementation, check if user exists or if token version matches
	// For now, generate new access token
	accessSecret := config.GetEnv("JWT_ACCESS_SECRET", "access_secret")
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	return accessToken.SignedString([]byte(accessSecret))
}

func (s *service) ForgotPassword(req ForgotPasswordRequest) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	code := utils.GenerateSixDigitCode()
	user.ResetCode = code
	user.ResetCodeExpiry = time.Now().Add(15 * time.Minute)

	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	utils.SendRecoveryEmail(user.Email, code)
	return nil
}

func (s *service) VerifyCode(req VerifyCodeRequest) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ResetCode != req.Code || time.Now().After(user.ResetCodeExpiry) {
		return errors.New("invalid or expired reset code")
	}

	return nil
}

func (s *service) ResetPassword(req ResetPasswordRequest) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ResetCode != req.Code || time.Now().After(user.ResetCodeExpiry) {
		return errors.New("invalid or expired reset code")
	}

	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashed
	user.ResetCode = ""
	user.ResetCodeExpiry = time.Time{}

	return s.repo.UpdateUser(user)
}

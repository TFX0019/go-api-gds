package auth

import (
	"errors"
	"time"

	"github.com/TFX0019/api-go-gds/features/plans"
	"github.com/TFX0019/api-go-gds/features/subscriptions"
	"github.com/TFX0019/api-go-gds/features/wallets"
	"github.com/TFX0019/api-go-gds/pkg/config"
	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, string, *UserResponse, error)
	VerifyEmail(token string) error
	RefreshToken(tokenString string) (string, error)
	ForgotPassword(req ForgotPasswordRequest) error
	VerifyCode(req VerifyCodeRequest) error
	ResetPassword(req ResetPasswordRequest) error
	UpdateAvatar(userID uint, avatarPath *string) (*UserResponse, error)
	UpdateName(userID uint, name string) (*UserResponse, error)
	VerifyAccount(req VerifyAccountRequest) (string, string, *UserResponse, error)
	ResendVerificationCode(req ResendCodeRequest) error
	ResendResetCode(req ForgotPasswordRequest) error
}

type service struct {
	repo      Repository
	plansRepo plans.Repository
}

func NewService(repo Repository, plansRepo plans.Repository) Service {
	return &service{repo: repo, plansRepo: plansRepo}
}

func (s *service) Register(req RegisterRequest) error {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		if existing.IsVerified {
			return errors.New("email already exists")
		}

		// Update existing unverified user
		if req.Password != req.ConfirmPassword {
			return errors.New("passwords do not match")
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		existing.Name = req.Name
		existing.Password = hashedPassword

		if err := s.repo.UpdateUser(existing); err != nil {
			return err
		}

		return s.generateAndSendCode(req.Email)
	}

	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		IsVerified: false,
		Wallet: wallets.Wallet{
			Balance:      30,
			LastRefillAt: time.Now(),
		},
		Subscription: subscriptions.Subscription{
			Status:    subscriptions.SubscriptionStatusExpired,
			ExpiresAt: time.Now(),
			ProductID: "free_tier",
		},
	}

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return s.generateAndSendCode(req.Email)
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

	// Logic to populate limits
	isPro := false
	planName := "Free Tier"
	maxCustomers := 20
	maxProducts := 20
	maxMaterials := 20
	maxTasks := 20

	if user.Subscription.ProductID != "" {
		plan, err := s.plansRepo.FindByProductID(user.Subscription.ProductID)
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	} else {
		// Fallback for empty product ID
		plan, err := s.plansRepo.FindByProductID("free_tier")
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	}

	if user.Subscription.Status == subscriptions.SubscriptionStatusActive && user.Subscription.ProductID != "free_tier" {
		isPro = true
		planName = user.Subscription.ProductID
	}

	userResponse := &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02 15:04:05"),
		IsPro:        isPro,
		Plan:         planName,
		MaxCustomers: maxCustomers,
		MaxProducts:  maxProducts,
		MaxMaterials: maxMaterials,
		MaxTasks:     maxTasks,
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

func (s *service) VerifyAccount(req VerifyAccountRequest) (string, string, *UserResponse, error) {
	// Find code
	vc, err := s.repo.FindVerificationCode(req.Email, req.Code)
	if err != nil {
		return "", "", nil, errors.New("invalid verification code")
	}

	if time.Now().After(vc.ExpiresAt) {
		return "", "", nil, errors.New("verification code expired")
	}

	// Find user
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}

	// Update user
	user.IsVerified = true
	if err := s.repo.UpdateUser(user); err != nil {
		return "", "", nil, err
	}

	// Delete code
	if err := s.repo.DeleteVerificationCode(req.Email); err != nil {
		// Just log error, don't fail flow
		// In production use a logger
	}

	// Login logic (generate tokens)
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	// Populate response (reuse logic from Login/UpdateAvatar or duplicate for now)
	isPro := false
	planName := "Free Tier"
	maxCustomers := 20
	maxProducts := 20
	maxMaterials := 20
	maxTasks := 20

	if user.Subscription.ProductID != "" {
		plan, err := s.plansRepo.FindByProductID(user.Subscription.ProductID)
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	} else {
		plan, err := s.plansRepo.FindByProductID("free_tier")
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	}

	if user.Subscription.Status == subscriptions.SubscriptionStatusActive && user.Subscription.ProductID != "free_tier" {
		isPro = true
		planName = user.Subscription.ProductID
	}

	userResponse := &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02 15:04:05"),
		IsPro:        isPro,
		Plan:         planName,
		MaxCustomers: maxCustomers,
		MaxProducts:  maxProducts,
		MaxMaterials: maxMaterials,
		MaxTasks:     maxTasks,
	}

	return accessToken, refreshToken, userResponse, nil
}

func (s *service) ResendVerificationCode(req ResendCodeRequest) error {
	// Find user
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		return errors.New("account already verified")
	}

	return s.generateAndSendCode(req.Email)
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
	user.ResetCodeExpiry = time.Now().Add(2 * time.Minute)

	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	utils.SendRecoveryCode(user.Email, code)
	return nil
}

func (s *service) ResendResetCode(req ForgotPasswordRequest) error {
	// Re-use logic
	return s.ForgotPassword(req)
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

func (s *service) UpdateAvatar(userID uint, avatarPath *string) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Avatar = avatarPath

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	// Logic to populate limits (duplicate from Login, should be refactored)
	isPro := false
	planName := "Free Tier"
	maxCustomers := 20
	maxProducts := 20
	maxMaterials := 20
	maxTasks := 20

	if user.Subscription.ProductID != "" {
		plan, err := s.plansRepo.FindByProductID(user.Subscription.ProductID)
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	} else {
		// Fallback for empty product ID
		plan, err := s.plansRepo.FindByProductID("free_tier")
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	}

	if user.Subscription.Status == subscriptions.SubscriptionStatusActive && user.Subscription.ProductID != "free_tier" {
		isPro = true
		planName = user.Subscription.ProductID
	}

	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Avatar:       user.Avatar, // Avatar is a pointer, if nil it stays nil in JSON
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02 15:04:05"),
		IsPro:        isPro,
		Plan:         planName,
		MaxCustomers: maxCustomers,
		MaxProducts:  maxProducts,
		MaxMaterials: maxMaterials,
		MaxTasks:     maxTasks,
	}, nil
}

func (s *service) UpdateName(userID uint, name string) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = name

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	// Logic to populate limits (duplicate from Login, should be refactored)
	isPro := false
	planName := "Free Tier"
	maxCustomers := 20
	maxProducts := 20
	maxMaterials := 20
	maxTasks := 20

	if user.Subscription.ProductID != "" {
		plan, err := s.plansRepo.FindByProductID(user.Subscription.ProductID)
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	} else {
		// Fallback for empty product ID
		plan, err := s.plansRepo.FindByProductID("free_tier")
		if err == nil && plan != nil {
			maxCustomers = plan.MaxCustomers
			maxProducts = plan.MaxProducts
			maxMaterials = plan.MaxMaterials
			maxTasks = plan.MaxTasks
		}
	}

	if user.Subscription.Status == subscriptions.SubscriptionStatusActive && user.Subscription.ProductID != "free_tier" {
		isPro = true
		planName = user.Subscription.ProductID
	}

	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02 15:04:05"),
		IsPro:        isPro,
		Plan:         planName,
		MaxCustomers: maxCustomers,
		MaxProducts:  maxProducts,
		MaxMaterials: maxMaterials,
		MaxTasks:     maxTasks,
	}, nil
}

func (s *service) generateAndSendCode(email string) error {
	// Delete existing codes
	if err := s.repo.DeleteVerificationCode(email); err != nil {
		// Log error but continue
	}

	// Generate new code
	code := utils.GenerateSixDigitCode()
	verificationCode := &VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}

	if err := s.repo.CreateVerificationCode(verificationCode); err != nil {
		return err
	}

	// Send email
	return utils.SendVerificationCode(email, code)
}

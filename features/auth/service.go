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
	Login(req LoginRequest, ip, userAgent string) (string, string, *UserResponse, error)
	VerifyEmail(token string) error
	RefreshToken(tokenString string) (string, error)
	ForgotPassword(req ForgotPasswordRequest) error
	VerifyCode(req VerifyCodeRequest) error
	ResetPassword(req ResetPasswordRequest) error
	UpdateAvatar(userID uint, avatarPath *string) (*UserResponse, error)
	UpdateName(userID uint, name string) (*UserResponse, error)
	VerifyAccount(req VerifyAccountRequest, ip, userAgent string) (string, string, *UserResponse, error)
	ResendVerificationCode(req ResendCodeRequest) error
	ResendResetCode(req ForgotPasswordRequest) error
	Logout(tokenString string) error
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

	// Find default role
	memberRole, err := s.repo.FindRoleByName("member")
	var roles []Role
	if err == nil && memberRole != nil {
		roles = append(roles, *memberRole)
	}
	// If role not found, we proceed without roles (or could handle error)

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
		Roles: roles,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return s.generateAndSendCode(req.Email)
}

func (s *service) Login(req LoginRequest, ip, userAgent string) (string, string, *UserResponse, error) {
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

	if !user.IsActive {
		return "", "", nil, errors.New("your account has been banned or deactivated")
	}

	return s.createSessionAndResponse(user, ip, userAgent)
}

func (s *service) createSessionAndResponse(user *User, ip, userAgent string) (string, string, *UserResponse, error) {
	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, roles)
	if err != nil {
		return "", "", nil, err
	}

	// Create Session
	session := &Session{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // Match refresh token expiry
		IPAddress: ip,
		UserAgent: userAgent,
		IsValid:   true,
	}

	if err := s.repo.CreateSession(session); err != nil {
		return "", "", nil, err
	}

	userResponse, err := s.buildUserResponse(user)
	if err != nil {
		return "", "", nil, err
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

func (s *service) VerifyAccount(req VerifyAccountRequest, ip, userAgent string) (string, string, *UserResponse, error) {
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
	}

	return s.createSessionAndResponse(user, ip, userAgent)
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
	// 1. Verify against DB Session
	session, err := s.repo.FindSessionByToken(tokenString)
	if err != nil {
		return "", errors.New("invalid or expired session")
	}

	if !session.IsValid || time.Now().After(session.ExpiresAt) {
		return "", errors.New("session expired or revoked")
	}

	// Fetch user to get current roles
	user, err := s.repo.FindByID(session.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	// Generate new access token
	accessSecret := config.GetEnv("JWT_ACCESS_SECRET", "access_secret")
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": session.UserID,
		"roles":   roles,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	return accessToken.SignedString([]byte(accessSecret))
}

func (s *service) Logout(tokenString string) error {
	return s.repo.RevokeSession(tokenString)
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

	return s.buildUserResponse(user)
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

	return s.buildUserResponse(user)
}

func (s *service) generateAndSendCode(email string) error {
	if err := s.repo.DeleteVerificationCode(email); err != nil {
		// Log error but continue
	}

	code := utils.GenerateSixDigitCode()
	verificationCode := &VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}

	if err := s.repo.CreateVerificationCode(verificationCode); err != nil {
		return err
	}

	return utils.SendVerificationCode(email, code)
}

func (s *service) buildUserResponse(user *User) (*UserResponse, error) {
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

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
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
		Roles:        roles,
	}, nil
}

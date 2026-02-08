package service

import (
	"context"
	"database/sql"
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
	"socialnet/internal/security"
	"strings"
	"unicode"
)

type AuthService struct {
	userRepo      *repository.UserRepository
	firebaseAuth  *security.FirebaseAuth
	initialAdmins []string
}

func NewAuthService(userRepo *repository.UserRepository, firebaseAuth *security.FirebaseAuth, initialAdmins []string) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		firebaseAuth:  firebaseAuth,
		initialAdmins: initialAdmins,
	}
}

func (s *AuthService) Register(reg *model.UserRegistration) (*model.User, error) {
	if err := security.ValidateEmail(reg.Email); err != nil {
		return nil, err
	}
	if err := security.ValidateUsername(reg.Username); err != nil {
		return nil, err
	}
	if err := s.ValidatePasswordStrength(reg.Password); err != nil {
		return nil, err
	}

	if _, err := s.userRepo.GetByEmail(reg.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	if _, err := s.userRepo.GetByUsername(reg.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := security.HashPassword(reg.Password)
	if err != nil {
		return nil, err
	}

	isAdmin := s.isInitialAdmin(reg.Email)

	user := &model.User{
		Email:        reg.Email,
		Username:     reg.Username,
		PasswordHash: hashedPassword,
		FullName:     reg.FullName,
		IsAdmin:      isAdmin,
	}

	id, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *AuthService) Login(login *model.UserLogin) (*model.User, error) {
	if err := security.ValidateEmail(login.Email); err != nil {
		return nil, errors.New("invalid email format")
	}

	user, err := s.userRepo.GetByEmail(login.Email)
	if err != nil {
		return nil, errors.New("account not found")
	}

	if !security.ComparePassword(user.PasswordHash, login.Password) {
		return nil, errors.New("incorrect password")
	}

	s.userRepo.UpdateOnlineStatus(user.ID, true)

	return user, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, idToken string) (*model.User, error) {
	if s.firebaseAuth == nil {
		return nil, errors.New("google auth not configured")
	}

	token, err := s.firebaseAuth.VerifyToken(ctx, idToken)
	if err != nil {
		return nil, errors.New("invalid google token")
	}

	user, err := s.userRepo.GetByFirebaseUID(token.UID)
	if err == nil {
		s.userRepo.UpdateOnlineStatus(user.ID, true)
		return user, nil
	}

	firebaseUser, err := s.firebaseAuth.GetUser(ctx, token.UID)
	if err != nil {
		return nil, errors.New("failed to get user info from google")
	}

	if existingUser, err := s.userRepo.GetByEmail(firebaseUser.Email); err == nil {
		s.userRepo.LinkFirebaseUID(existingUser.ID, token.UID)
		s.userRepo.UpdateOnlineStatus(existingUser.ID, true)
		return existingUser, nil
	}

	username := s.generateUsername(firebaseUser.Email)
	isAdmin := s.isInitialAdmin(firebaseUser.Email)

	newUser := &model.User{
		Email:       firebaseUser.Email,
		Username:    username,
		FullName:    firebaseUser.DisplayName,
		AvatarURL:   firebaseUser.PhotoURL,
		FirebaseUID: sql.NullString{String: token.UID, Valid: true},
		IsAdmin:     isAdmin,
	}

	id, err := s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	newUser.ID = id
	s.userRepo.UpdateOnlineStatus(newUser.ID, true)
	return newUser, nil
}

func (s *AuthService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func (s *AuthService) GetPasswordRequirements() []string {
	return []string{
		"At least 8 characters",
		"At least one uppercase letter (A-Z)",
		"At least one lowercase letter (a-z)",
		"At least one number (0-9)",
		"At least one special character (!@#$%^&*)",
	}
}

func (s *AuthService) isInitialAdmin(email string) bool {
	for _, adminEmail := range s.initialAdmins {
		if strings.EqualFold(adminEmail, email) {
			return true
		}
	}
	return false
}

func (s *AuthService) generateUsername(email string) string {
	parts := strings.Split(email, "@")
	base := strings.ReplaceAll(parts[0], ".", "_")

	username := base
	counter := 1

	for {
		if _, err := s.userRepo.GetByUsername(username); err != nil {
			return username
		}
		counter++
		username = base + string(rune('0'+counter))
		if counter > 99 {
			return base + "_" + strings.ReplaceAll(parts[1], ".", "")
		}
	}
}

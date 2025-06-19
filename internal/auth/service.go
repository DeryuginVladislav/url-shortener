package auth

import (
	"errors"
	"url_shortener/internal/user"
	"url_shortener/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserReposytory
}

func NewAuthService(userRepository di.IUserReposytory) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (s *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := s.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = s.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
func (s *AuthService) Login(email, password string) (string, error) {
	u, _ := s.UserRepository.FindByEmail(email)
	if u == nil {
		return "", errors.New(ErrWrongCredentials)
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}

	return u.Email, nil
}

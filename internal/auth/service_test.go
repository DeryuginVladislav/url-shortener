package auth_test

import (
	"testing"
	"url_shortener/internal/auth"
	"url_shortener/internal/user"
)

type MockUserReposytory struct{}

func (m *MockUserReposytory) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "lol@gmail.com",
	}, nil
}
func (m *MockUserReposytory) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "lol@gmail.com"
	authService := auth.NewAuthService(&MockUserReposytory{})
	email, err := authService.Register(initialEmail, "2", "vlad")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatalf("Email %s do not match %s", email, initialEmail)
	}

}

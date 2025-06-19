package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/configs"
	"url_shortener/internal/auth"
	"url_shortener/internal/user"
	"url_shortener/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{DB: gormDB})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}
func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("afa4@gmail.com", "$2a$10$kJbhcWp5SHeOFLHONZTRZ.4F7S6Qabbx.bqA1Uk5JHky.TyjdfZK6")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "afa4@gmail.com",
		Password: "asdasd",
	})

	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)

	handler.Login()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", 200, w.Code)
	}

}
func TestRegisterSuccess2(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "afa5@gmail.com",
		Password: "asdasd",
		Name:     "vlad",
	})

	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)

	handler.Register()(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", 201, w.Code)
	}

}

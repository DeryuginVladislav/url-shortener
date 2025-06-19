package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"
	"url_shortener/internal/auth"
	"url_shortener/internal/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "afa4@gmail.ru",
		Password: "$2a$10$kJbhcWp5SHeOFLHONZTRZ.4F7S6Qabbx.bqA1Uk5JHky.TyjdfZK6",
		Name:     "Vasya",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email=?", "afa4@gmail.ru").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	//Prepere
	db := initDB()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "afa4@gmail.ru",
		Password: "asdasd",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	var loginResp auth.LoginResponse
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		t.Fatal("Expected token")
	}

}
func TestLoginFail(t *testing.T) {
	//Prepere
	db := initDB()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())

	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "afa4@gmail.ru",
		Password: "asdas",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 401 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}

}

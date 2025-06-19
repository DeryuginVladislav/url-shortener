package jwt_test

import (
	"testing"
	"url_shortener/pkg/jwt"
)

func TestJWTCreate(t *testing.T) {
	const email = "afa4@gmail.com"
	jwtService := jwt.NewJWT("bE/QxHUjaqkUK5E1BxpPmXBBawrc3xG+I6NnOumRQuA=")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is invalid")
	}
	if data.Email != email {
		t.Fatalf("Email %s not equal %s", data.Email, email)
	}
}

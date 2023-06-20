package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return []byte{}, err
	}
	return hashedPass, nil
}

func CheckPasswordHash(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func CreateJwt(id int, secret string, tokenType string) (string, error) {
	exp := time.Hour
	if tokenType == "refresh" {
		exp = 60 * (24 * time.Hour)
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("chirpy-%s", tokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		Subject:   fmt.Sprintf("%v", id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJwt(token string, secret string, tokenType string) (string, error) {
	type myCustomClaims struct {
		jwt.RegisteredClaims
	}
	jwtToken, err := jwt.ParseWithClaims(token, &myCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	if !jwtToken.Valid {
		return "", errors.New("Invalid token")
	}

	claims, ok := jwtToken.Claims.(*myCustomClaims)
	if !ok {
		return "", errors.New("Something went wrong")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return "", errors.New("Token expired")
	}

	if claims.Issuer != fmt.Sprintf("chirpy-%s", tokenType) {
		return "", errors.New("Wrong token type")
	}

	id, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetBearerToken(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Missing token")
	}
	token := authHeader[7:]

	return token, nil
}

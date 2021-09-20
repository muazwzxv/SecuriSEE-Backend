package service

import (
	"Oracle-Hackathon-BE/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtWrapper struct {
	SecretKey    string
	Issuer       string
	ExpiredHours int64
}

type JwtClaims struct {
	IC   string
	ID   string
	Role []string
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(user *config.UserJwt) (string, error) {
	claims := &JwtClaims{
		ID:   user.ID,
		IC:   user.IC,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpiredHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if signed, err := token.SignedString([]byte(j.SecretKey)); err != nil {
		return "", err
	} else {
		return signed, nil
	}
}

func (j *JwtWrapper) VerifyToken(token string) (jwt.MapClaims, error) {
	if token != "" {
		return nil, errors.New("token is null")
	}

	validate, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetInstance().GetJWTSecret()), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := validate.Claims.(jwt.MapClaims); !ok && !validate.Valid {
		return nil, err
	} else {
		return claims, nil
	}
}

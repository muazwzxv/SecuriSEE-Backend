package service

import (
	"Oracle-Hackathon-BE/config"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtWrapper struct {
	SecretKey    string
	Issuer       string
	ExpiredHours int64
}

type JwtClaims struct {
	IC   string
	ID   uint64
	Role []string
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(user *config.UserJwt) (string, error) {
	claims := &JwtClaims{
		ID:   user.ID,
		IC:   user.Ic,
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

func (j *JwtWrapper) VerifyToken(token string) error {
	if token != "" {
		return errors.New("token is null")
	}

	validate, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.CFG.GetJWTSecret()), nil
	})
	if err != nil {
		return err
	}

	if _, ok := validate.Claims.(jwt.MapClaims); !ok && !validate.Valid {
		return err
	}
	return nil
}

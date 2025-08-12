package utils

import (
	"lms-bootcamp/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Utils struct {
	JwtSecret string
}

func NewUtils() *Utils {
	conf, err := config.NewConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	if conf == nil {
		panic("Config is nil")
	}
	return &Utils{
		JwtSecret: conf.JWTSecret,
	}
}


func (u *Utils) GenerateJWTToken(data interface{}) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.JwtSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (u *Utils) ValidateJWTToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func (token *jwt.Token) (interface{}, error) {
		return []byte(u.JwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}
package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func GenerateVerificationCode(salt string) string {
	b := []byte(salt + time.Now().String())
	return fmt.Sprintf("%x", md5.Sum(b))
}

func ReplacePlaceholders(message string, params map[string]interface{}) string {
	if len(message) == 0 {
		return ""
	}
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}

func CreateToken(alg, secretKey string, userId int, exp time.Duration) string {
	method := jwt.GetSigningMethod(alg)
	token := jwt.NewWithClaims(method, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(exp).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func ParseToken(secretKey, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	userId := token.Claims.(jwt.MapClaims)["user_id"].(float64)
	return int(userId), nil
}

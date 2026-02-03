package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 1. Struct giữ dependency
type JWTWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// 2. Constructor: Chỉ nhận tham số, KHÔNG tự load config
func NewJWTWrapper(secretKey, issuer string, expirationHours int64) *JWTWrapper {
	// Set default nếu cần thiết (hoặc để layer ngoài lo)
	if issuer == "" {
		issuer = "my-app-backend"
	}

	return &JWTWrapper{
		SecretKey:       secretKey,
		Issuer:          issuer,
		ExpirationHours: expirationHours,
	}
}

// 3. Custom Claims
type MyCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 4. GenerateJWT (Dùng w.SecretKey và w.ExpirationHours)
func (w *JWTWrapper) GenerateJWT(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(w.ExpirationHours) * time.Hour)

	claims := MyCustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    w.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Lưu ý: w.SecretKey là string, cần ép kiểu về []byte
	return token.SignedString([]byte(w.SecretKey))
}

// 5. ValidateJWT
func (w *JWTWrapper) ValidateJWT(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(w.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

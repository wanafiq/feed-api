package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	BearerPrefix = "Bearer "
)

type CustomClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func Hash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyHash(hashedValue, value string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	return err == nil
}

func GenerateJWT(user *models.User, secret string, expiredAt time.Time, issuer string, audience string) (string, error) {
	secretKey := []byte(secret)
	now := time.Now()

	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID,
			ExpiresAt: expiredAt.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    issuer,
			Audience:  audience,
		},
		Role: user.Role.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func ParseAndValidateJWT(authHeader, secret string) (*CustomClaims, error) {
	parts := strings.Fields(authHeader)

	if authHeader == "" || !strings.HasPrefix(authHeader, BearerPrefix) || len(parts) != 2 {
		return nil, constants.ErrInvalidAuthHeader
	}

	tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, constants.ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims.ExpiresAt < time.Now().Unix() {
		return nil, constants.ErrExpiredJWT
	}

	return claims, nil
}

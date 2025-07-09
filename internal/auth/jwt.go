package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type JWTService struct {
	accessSecret  string
	refreshSecret string
}

func NewJWTService(accessSecret string, refreshSecret string) *JWTService {
	return &JWTService{accessSecret: accessSecret, refreshSecret: refreshSecret}
}

func (s *JWTService) GenerateTokens(userId int64) (string, string, error) {
	accessToken, err := s.createAccessToken(userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.createRefreshToken(userId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *JWTService) GetUserIdFromToken(token *jwt.Token) (int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("could not parse claims")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("could not get subject from claims: %w", err)
	}

	userId, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parser user ID from subject: %w", err)
	}

	return userId, nil

}
func (s *JWTService) createAccessToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userId),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iss": "orders-service",
		"aud": "orders-service-users",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

func (s *JWTService) createRefreshToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userId),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

func (s *JWTService) validateToken(tokenString string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	return s.validateToken(tokenString, []byte(s.accessSecret))
}
func (s *JWTService) ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return s.validateToken(tokenString, []byte(s.refreshSecret))
}

package bearerrepo

import (
	"errors"
	"fmt"
	"time"

	"mizito/internal/database"
	userdto "mizito/pkg/models/dtos/user"

	"github.com/golang-jwt/jwt/v4"
)

type BearerRepository interface {
	AuthorizeBearerUser(tokenString string) (*userdto.UserClaims, error)
	GenerateTokens(userID uint, username string) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
	BlacklistToken(tokenString string, ttl time.Duration) error
	IsTokenBlacklisted(tokenString string) (bool, error)
}

type jwtRepository struct {
	secret string
	redis  *database.RedisHandler
}

func NewJwtRepository(secret string, redis *database.RedisHandler) BearerRepository {
	return &jwtRepository{
		secret: secret,
		redis:  redis,
	}
}

func (jr *jwtRepository) AuthorizeBearerUser(tokenString string) (*userdto.UserClaims, error) {

	isBlacklisted, err := jr.IsTokenBlacklisted(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if isBlacklisted {
		return nil, errors.New("token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenString, &userdto.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jr.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*userdto.UserClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (jr *jwtRepository) GenerateTokens(userID uint, username string) (string, string, error) {
	accessClaims := userdto.UserClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "mizito",
		},
	}
	refreshClaims := userdto.UserClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "mizito",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessString, err := accessToken.SignedString([]byte(jr.secret))
	if err != nil {
		return "", "", err
	}
	refreshString, err := refreshToken.SignedString([]byte(jr.secret))
	if err != nil {
		return "", "", err
	}
	return accessString, refreshString, nil
}

func (jr *jwtRepository) RefreshTokens(refreshToken string) (string, string, error) {

	isBlacklisted, err := jr.IsTokenBlacklisted(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if isBlacklisted {
		return "", "", errors.New("refresh token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(refreshToken, &userdto.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jr.secret), nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := token.Claims.(*userdto.UserClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	return jr.GenerateTokens(claims.UserID, claims.Username)
}

func (jr *jwtRepository) BlacklistToken(tokenString string, ttl time.Duration) error {
	return jr.redis.SetBlacklistedToken(tokenString, ttl)
}

func (jr *jwtRepository) IsTokenBlacklisted(tokenString string) (bool, error) {
	return jr.redis.IsTokenBlacklisted(tokenString)
}

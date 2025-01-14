package user_dto

import "github.com/golang-jwt/jwt/v4"

type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

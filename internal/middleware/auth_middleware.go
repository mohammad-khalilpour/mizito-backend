package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"mizito/internal/env"
	"strings"
)

func AuthMiddleware(c *fiber.Ctx) error {
	if strings.HasPrefix(c.Path(), "/api/auth") || strings.HasPrefix(c.Path(), "/user") {
		return c.Next()
	}

	token := c.Get("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := ValidateTokenAndExtractUserID(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: " + err.Error(),
		})
	}

	c.Locals("userID", userID)
	return c.Next()
}

// CustomClaims defines the claims stored in the JWT token
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateTokenAndExtractUserID validates the token and extracts the user ID
func ValidateTokenAndExtractUserID(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, errors.New("missing token")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.Config{}.AuthorizationSecret), nil
	})
	if err != nil {
		return 0, errors.New("invalid token")
	}

	// Validate the token and extract claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("unauthorized")
}

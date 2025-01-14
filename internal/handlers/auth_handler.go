package handlers

import (
	"fmt"
	basicrepo "mizito/internal/repositories/auth/basic"
	bearerrepo "mizito/internal/repositories/auth/bearer"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
	Authorize(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type authHandler struct {
	jwtRepo   bearerrepo.BearerRepository
	basicRepo basicrepo.BasicRepository
}

func NewAuthHandler(jwtRepo bearerrepo.BearerRepository, basicRepo basicrepo.BasicRepository) AuthHandler {
	return &authHandler{
		jwtRepo:   jwtRepo,
		basicRepo: basicRepo,
	}
}

func (ah *authHandler) Login(ctx *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&credentials); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	authenticated, userID, err := ah.basicRepo.AuthenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}
	if !authenticated {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	token, refreshToken, err := ah.jwtRepo.GenerateTokens(userID, credentials.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate tokens",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func (ah *authHandler) Refresh(ctx *fiber.Ctx) error {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if payload.RefreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}
	accessToken, newRefreshToken, err := ah.jwtRepo.RefreshTokens(payload.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func (ah *authHandler) Authorize(ctx *fiber.Ctx) error {

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}

	var tokenString string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Authorization header format",
		})
	}

	claims, err := ah.jwtRepo.AuthorizeBearerUser(tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx.Locals("userID", claims.UserID)
	ctx.Locals("username", claims.Username)

	return ctx.Next()
}

func (ah *authHandler) Logout(ctx *fiber.Ctx) error {

	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if payload.RefreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	claims, err := ah.jwtRepo.AuthorizeBearerUser(payload.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	expiration := time.Unix(claims.ExpiresAt, 0)
	ttl := time.Until(expiration)
	if ttl <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token has already expired",
		})
	}

	err = ah.jwtRepo.BlacklistToken(payload.RefreshToken, ttl)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to blacklist refresh token",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
	basichandler "mizito/internal/repositories/auth/basic"
	bearerhandler "mizito/internal/repositories/auth/bearer"
)

func InitAuth(r *Router, secret string, redis *database.RedisHandler, db *database.DatabaseHandler) {
	jwtRepo := bearerhandler.NewJwtRepository(secret, redis)
	basicRepo := basichandler.NewBasicHandler(db)

	authHandler := handlers.NewAuthHandler(jwtRepo, basicRepo)

	authGroup := r.App.Group("/api/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.Refresh)
}

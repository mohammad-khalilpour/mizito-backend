package router

import (
	"fmt"
	"mizito/internal/database"
	"mizito/internal/env"
	"mizito/internal/websocket/websocket"
)
import websocketfiber "github.com/gofiber/contrib/websocket"

func InitSocket(r *Router, redis *database.RedisHandler, mongo *database.MongoHandler, postgreSql *database.DatabaseHandler, env *env.Config) {

	fmt.Println("initializing socket routes...")

	socketManager := websocket.NewChannelHandler(redis, mongo, postgreSql, env)

	r.App.Get("/ws/:id", websocketfiber.New(socketManager.Register))
}

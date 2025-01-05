package router

import (
	"fmt"
	"mizito/internal/websocket/websocket"
)
import websocketfiber "github.com/gofiber/contrib/websocket"

func InitSocket(r *Router) {

	fmt.Println("initializing socket routes...")

	ch := websocket.NewChannelHandler()

	r.App.Get("/ws/:id", websocketfiber.New(ch.Register))
}

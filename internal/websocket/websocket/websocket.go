package websocket

import "github.com/gofiber/contrib/websocket"


type WebSocketManager interface {
	AddSocket(id int, conn *websocket.Conn)
	RemoveSocket(id int) error
	GetSocketByID(id int) (*websocket.Conn, error)
}
package websocket

import "github.com/gofiber/contrib/websocket"

type WebSocketManager interface {
	AddSocket(id uint, conn *websocket.Conn)
	RemoveSocket(id uint) error
	GetSocketByID(id uint) (*websocket.Conn, error)
}

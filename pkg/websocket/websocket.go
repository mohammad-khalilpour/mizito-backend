package websocket

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketManager interface {
	AddSocket(id string, conn *websocket.Conn)
	RemoveSocket(id string) error
	GetSocketByID(id string) (*websocket.Conn, error)
}


type NotificationRouter interface {
	send(ch <-chan []byte, ids []int)
}


type ChannelManager interface {
	WebSocketManager
	NotificationRouter
}


type channelManager struct {
	sockets map[string]*websocket.Conn
	mu sync.Mutex
}


func NewChannelManager() ChannelManager{
	return &channelManager{}
}

func (m *channelManager) AddSocket(id string, conn *websocket.Conn){
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sockets[id] = conn
}

func (m *channelManager) RemoveSocket(id string) error{
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sockets[id]; !ok {
		return fmt.Errorf("socket with id %s not found", id)
	}


	delete(m.sockets, id)

	return nil
}


func (m *channelManager) GetSocketByID(id string) (*websocket.Conn, error){
	sck, ok := m.sockets[id]
	if !ok {
		return nil, fmt.Errorf("no such key %s found", id)
	}
	return sck, nil
}


func (m *channelManager) send(ch <-chan []byte, ids []int) {
	//TODO
}
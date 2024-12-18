package websocket

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketManager interface {
	AddSocket(id int, conn *websocket.Conn)
	RemoveSocket(id int) error
	GetSocketByID(id int) (*websocket.Conn, error)
}


type NotificationRouter interface {
	send(ch <-chan *ChannelMessage)
}

type ChannelMessage struct {
	Notification *NotificationMessage
	Ids []int
}


type NotificationMessage struct {
	Payload []byte
	EventType string
}


type ChannelManager interface {
	WebSocketManager
	NotificationRouter
}

type channelManager struct {
	sockets map[int]*websocket.Conn
	mu sync.Mutex
}


func NewChannelManager(size int) ChannelManager{


	return &channelManager{
		
	}
}

func (m *channelManager) AddSocket(id int, conn *websocket.Conn){
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sockets[id] = conn
}

func (m *channelManager) RemoveSocket(id int) error{
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sockets[id]; !ok {
		return fmt.Errorf("socket with id %s not found", id)
	}


	delete(m.sockets, id)

	return nil
}


func (m *channelManager) GetSocketByID(id int) (*websocket.Conn, error){
	sck, ok := m.sockets[id]
	if !ok {
		return nil, fmt.Errorf("no such key %d found", id)
	}
	return sck, nil
}


func (m *channelManager) send(ch <-chan *ChannelMessage) {

	for msg := range ch {
		for id := range msg.Ids {
			m.sendMessage(*msg.Notification, id)
		}
	}
}

func (m *channelManager) sendMessage(msg NotificationMessage, id int) {
	if conn, err := m.GetSocketByID(id); err != nil {
		// log error and continue
	} else {
		if err := conn.WriteJSON(msg); err != nil {
			// log error and continue
		}
	}
}
package websocket

import (
	"fmt"
	"mizito/pkg/models/dtos"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type ChannelManager interface {
	WebSocketManager
	EventRouter
}

type channelManager struct {
	sockets sync.Map
}

func NewChannelManager(size int) ChannelManager {

	return &channelManager{}
}

func (m *channelManager) AddSocket(id int, conn *websocket.Conn) {
	m.sockets.Store(id, conn)
}

func (m *channelManager) RemoveSocket(id int) error {

	if _, ok := m.sockets.Load(id); !ok {
		return fmt.Errorf("socket with id %d not found", id)
	}

	m.sockets.Delete(id)

	return nil
}

func (m *channelManager) GetSocketByID(id int) (*websocket.Conn, error) {
	sck, ok := m.sockets.Load(id)
	if !ok {
		return nil, fmt.Errorf("no such key %d found", id)
	}
	conn, ok := sck.(*websocket.Conn)
	if !ok {
		return nil, fmt.Errorf("failed to parse value into websocket connection")
	}
	return conn, nil
}

func (m *channelManager) Publish(ch <-chan *dtos.EventMessage) {

	for msg := range ch {
		for id := range msg.Ids {
			m.publishEvent(msg.Event, id)
		}
	}
}

func (m *channelManager) publishEvent(msg *dtos.Event, id int) {
	if conn, err := m.GetSocketByID(id); err != nil {
		// log error and continue
	} else {
		if err := conn.WriteJSON(msg); err != nil {
			// log error and continue
		}
	}
}

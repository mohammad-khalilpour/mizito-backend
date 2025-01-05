package websocket

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"mizito/pkg/models/dtos"
)

type SocketManager interface {
	WebSocketManager
	EventRouter
}

type socketManager struct {
	sockets   map[uint]*websocket.Conn
	eventChan chan *dtos.EventMessage
}

func NewSocketHandler() SocketManager {

	ch := make(chan *dtos.EventMessage)

	return &socketManager{
		sockets:   make(map[uint]*websocket.Conn),
		eventChan: ch,
	}
}

func (m *socketManager) AddSocket(id uint, conn *websocket.Conn) {
	m.sockets[id] = conn
}

func (m *socketManager) RemoveSocket(id uint) error {

	if _, ok := m.sockets[id]; !ok {
		return fmt.Errorf("socket with id %d not found", id)
	}

	delete(m.sockets, id)

	return nil
}

func (m *socketManager) GetSocketByID(id uint) (*websocket.Conn, error) {

	sck, ok := m.sockets[id]

	if !ok {
		return nil, fmt.Errorf("no such key %d found", id)
	}

	return sck, nil
}

func (m *socketManager) AddToPublishChan(e *dtos.EventMessage) {
	m.eventChan <- e
}

func (m *socketManager) Publish() {

	for msg := range m.eventChan {
		for _, id := range msg.Ids {
			m.publishEvent(msg.Event, id)
		}
	}
}

func (m *socketManager) publishEvent(msg *dtos.Event, id uint) {
	if conn, err := m.GetSocketByID(id); err != nil {
		// log error and continue
	} else {
		if err := conn.WriteJSON(msg); err != nil {
			// log error and continue
		}
	}
}

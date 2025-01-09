package websocket

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"mizito/pkg/models/dtos"
)

type EventRouter interface {
	SendEvent(e *dtos.WebSocketMessage)
}

type WebSocketManager interface {
	AddSocket(id uint, conn *websocket.Conn)
	RemoveSocket(id uint) error
	GetSocketByID(id uint) (*websocket.Conn, error)
}

type SocketManager interface {
	WebSocketManager
	EventRouter
}

type socketManager struct {
	sockets   map[uint]*websocket.Conn
	eventChan chan *dtos.WebSocketMessage
}

func NewSocketHandler() SocketManager {

	ch := make(chan *dtos.WebSocketMessage)

	sm := &socketManager{
		sockets:   make(map[uint]*websocket.Conn),
		eventChan: ch,
	}
	go sm.publish()

	return sm
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

func (m *socketManager) SendEvent(e *dtos.WebSocketMessage) {
	m.eventChan <- e
}

func (m *socketManager) publish() {

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

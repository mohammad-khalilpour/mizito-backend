package websocket

import (
	"mizito/pkg/models/dtos"
	"strconv"
)

import "github.com/gofiber/contrib/websocket"

type ChannelHandler struct {
	socketManager ChannelManager

	// all the events received from socket is sent through this channel
	socketChan chan dtos.EventMessage
	// used to send events to users we want to broadcast message to
	// messages that are supposed to broadcast will pass through this
	eventChan chan<- dtos.EventMessage
	// channel for storing messages inside db
	// events published inside the `socketChan` are published to `eventChan` and `messageChan` based on their type
	// events that are of type message need to be sent to both `eventChan` and `messageChan`
	messageChan chan dtos.EventMessage
}

func (chm ChannelHandler) HandleEvent() {
	for event := range chm.socketChan {
		switch event.Event.EventType {
		case dtos.Message:
			chm.messageChan <- event
			chm.eventChan <- event
		case dtos.Notification:
			chm.eventChan <- event
		}
	}
}

func (chm ChannelHandler) Register(c *websocket.Conn) {
	sid := c.Params("id")
	// middleware checks id being integer
	id, _ := strconv.ParseInt(sid, 10, 32)

	chm.socketManager.AddSocket(int(id), c)
	defer chm.socketManager.RemoveSocket(int(id))

	for {

		var e dtos.EventMessage
		if err := c.ReadJSON(&e); err != nil {
			return
		}
		chm.eventChan <- e

	}

}

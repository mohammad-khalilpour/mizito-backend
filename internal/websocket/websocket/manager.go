package websocket

import (
	"encoding/json"
	"mizito/internal/database"
	"mizito/pkg/models/dtos"
	"strconv"
)

import "github.com/gofiber/contrib/websocket"

type ChannelHandler struct {
	socketManager SocketManager
	RedisClient   database.RedisHandler
	MongoClient   database.MongoHandler
}

func NewChannelHandler() *ChannelHandler {

	sm := NewSocketHandler()
	go sm.Publish()

	return &ChannelHandler{
		socketManager: sm,
	}
}

func (chm ChannelHandler) ProcessEvents() {
	for e := range chm.RedisClient.GetSubscribeChan() {
		var event dtos.EventMessage
		if err := json.Unmarshal(e, &event); err != nil {
			continue
		}
		chm.socketManager.AddToPublishChan(&event)
	}
}

func (chm ChannelHandler) Register(c *websocket.Conn) {
	sid := c.Params("id")
	// middleware checks id being integer
	id, _ := strconv.ParseInt(sid, 10, 32)

	chm.socketManager.AddSocket(uint(id), c)

	for {
		var (
			e    dtos.EventMessage
			eRaw []byte
			err  error
		)
		if _, eRaw, err = c.ReadMessage(); err != nil {
			//handle related error
			continue
		}
		if err := json.Unmarshal(eRaw, &e); err != nil {
			continue
		}
		chm.RedisClient.AddToPublishChan(eRaw)
		if e.Event.EventType == dtos.Message {
			chm.MongoClient.StoreMessage(eRaw)
		}

	}

}

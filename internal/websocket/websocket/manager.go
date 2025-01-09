package websocket

import (
	"encoding/json"
	"mizito/internal/database"
	"mizito/internal/repositories"
	"mizito/pkg/models/dtos"
	message_dto "mizito/pkg/models/dtos/message"
	"strconv"
)

import "github.com/gofiber/contrib/websocket"

type ChannelHandler struct {
	socketManager SocketManager
	RedisClient   database.RedisHandler
	MongoClient   database.MongoHandler
	ProjectDetail repositories.ProjectDetail
}

func NewChannelHandler(redisClient database.RedisHandler, mongoClient database.MongoHandler) *ChannelHandler {

	sm := NewSocketHandler()
	go sm.Publish()

	chHandler := &ChannelHandler{
		socketManager: sm,
		RedisClient:   redisClient,
		MongoClient:   mongoClient,
	}

	go chHandler.ProcessEvents()

	return chHandler
}

func (chm ChannelHandler) ProcessEvents() {
	for e := range chm.RedisClient.GetSubscribeChan() {
		var event dtos.EventMessage
		if err := json.Unmarshal(e, &event); err != nil {
			continue
		}
		if event.Event.EventType == dtos.Message {
			chm.processMsg(&event)
		}

		chm.socketManager.AddToPublishChan(&event)
	}
}

func (chm ChannelHandler) processMsg(event *dtos.EventMessage) {
	var msg message_dto.Message
	if err := json.Unmarshal([]byte(event.Event.Payload), &msg); err != nil {
		// log for error or produce to kafka error queue
	}
	if members, err := chm.ProjectDetail.GetProjectMembers(msg.Project); err != nil {
		// log again for errors
	} else {
		for _, member := range members {
			event.Ids = append(event.Ids, member.User.ID)
		}
	}
}

func (chm ChannelHandler) Register(c *websocket.Conn) {
	sid := c.Params("id")
	// middleware checks id being integer
	id, _ := strconv.ParseInt(sid, 10, 32)

	chm.socketManager.AddSocket(uint(id), c)

	for {
		var (
			e    dtos.Event
			eRaw []byte
			err  error
		)
		if err = c.ReadJSON(&e); err != nil {
			//handle related error
			continue
		}
		if eRaw, err = json.Marshal(e.Payload); err != nil {
			continue
		}
		if e.EventType == dtos.Message {
			go chm.MongoClient.StoreMessage(eRaw)

		}

	}

}

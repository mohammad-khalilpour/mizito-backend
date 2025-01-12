package websocket

import (
	"encoding/json"
	"fmt"
	"mizito/internal/database"
	"mizito/internal/env"
	"mizito/internal/repositories"
	"mizito/pkg/models/dtos"
	messagedto "mizito/pkg/models/dtos/message"
	"strconv"
	"time"
)

import "github.com/gofiber/contrib/websocket"

type ChannelRepository struct {
	socketManager SocketManager
	messageRepo   repositories.MessageChannelRepository
	ProjectDetail repositories.ProjectDetailRepo
}

func NewChannelHandler(redis *database.RedisHandler, mongo *database.MongoHandler, postgreSql *database.DatabaseHandler, env *env.Config) *ChannelRepository {

	sm := NewSocketHandler()

	chHandler := &ChannelRepository{
		socketManager: sm,
		messageRepo:   repositories.NewMessageRepository(redis, mongo, env),
		ProjectDetail: repositories.NewProjectRepository(postgreSql),
	}

	go chHandler.ProcessEvents()

	return chHandler
}

func (chm ChannelRepository) ProcessEvents() {
	for e := range chm.messageRepo.SubscribeEvent() {
		var event dtos.WebSocketMessage
		event.Event = &e
		if event.Event.EventType == dtos.Message {
			chm.processMsg(&event)
		}
		fmt.Println("final event is : ", event)
		chm.socketManager.SendEvent(&event)
	}
}

func (chm ChannelRepository) processMsg(event *dtos.WebSocketMessage) {
	var msg messagedto.Message
	fmt.Println(event.Event.Payload)
	if members, err := chm.ProjectDetail.GetProjectMembers(msg.Project); err != nil {
		// log again for errors
	} else {
		for _, member := range members {
			event.Ids = append(event.Ids, member.UserID)
		}
	}
}

func (chm ChannelRepository) Register(c *websocket.Conn) {
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
			fmt.Println(err.Error())
			continue
		}
		e.Payload.CreatedAt = time.Now()
		if eRaw, err = json.Marshal(e); err != nil {
			fmt.Println(err.Error())
			continue
		}

		if e.EventType == dtos.Message {
			go chm.messageRepo.PublishMsg(eRaw)

		}

	}

}

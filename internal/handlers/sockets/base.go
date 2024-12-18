package sockets

import (
	wsmanager "mizito/pkg/websocket"
	"strconv"

	"github.com/gofiber/contrib/websocket"
)



type ChannelHandler struct {
	socketManager wsmanager.ChannelManager
}


type Event struct {

}



func (chm ChannelHandler) Register (c *websocket.Conn) {
	sid := c.Params("id")
	// middleware checks id being integer
	id, _ := strconv.ParseInt(sid, 10, 32)

	chm.socketManager.AddSocket(int(id), c)
	defer chm.socketManager.RemoveSocket(int(id))


	for {

		var e Event
		if err := c.ReadJSON(&e); err != nil {
			return
		}
		//DO THE REST


	}



}
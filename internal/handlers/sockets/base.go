package sockets

import (
	wsmanager "mizito/pkg/websocket"

	"github.com/gofiber/contrib/websocket"
)



type ChannelHandler struct {
	socketManager wsmanager.ChannelManager
}


type Event struct {

}



func (chm ChannelHandler) Register (c *websocket.Conn) {
	id := c.Params("id")

	chm.socketManager.AddSocket(id, c)
	defer chm.socketManager.RemoveSocket(id)


	for {

		var e Event
		if err := c.ReadJSON(&e); err != nil {
			return
		}
		//DO THE REST


	}



}
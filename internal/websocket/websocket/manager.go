package websocket

import (
	"strconv"	
)

import "github.com/gofiber/contrib/websocket"



type ChannelHandler struct {
	socketManager ChannelManager
	eventChan chan<- ChannelMessage
}




func (chm ChannelHandler) Register (c *websocket.Conn) {
	sid := c.Params("id")
	// middleware checks id being integer
	id, _ := strconv.ParseInt(sid, 10, 32)

	chm.socketManager.AddSocket(int(id), c)
	defer chm.socketManager.RemoveSocket(int(id))


	for {

		var e ChannelMessage
		if err := c.ReadJSON(&e); err != nil {
			return
		}
		chm.eventChan <- e
		
	}



}
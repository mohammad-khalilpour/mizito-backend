package websocket

import "mizito/pkg/models/dtos"

type EventRouter interface {
	Publish(ch <-chan *dtos.EventMessage)
}

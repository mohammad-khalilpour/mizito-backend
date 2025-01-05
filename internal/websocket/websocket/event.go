package websocket

import "mizito/pkg/models/dtos"

type EventRouter interface {
	Publish()
	AddToPublishChan(e *dtos.EventMessage)
}

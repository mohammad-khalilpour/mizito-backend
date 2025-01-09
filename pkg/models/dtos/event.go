package dtos

import message_dto "mizito/pkg/models/dtos/message"

type EventType string

const (
	Message      EventType = "message"
	Notification EventType = "notification"
)

type Event struct {
	// for later use, struct must get generic attribute for Payload
	Payload   message_dto.Message `validate:"dive" json:"payload"`
	EventType EventType           `validate:"oneof='message notification'" json:"event_type"`
}

type EventMessage struct {
	Event *Event
	Ids   []uint
}

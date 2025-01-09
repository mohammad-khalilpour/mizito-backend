package dtos

type EventType string

const (
	Message      EventType = "message"
	Notification EventType = "notification"
)

type Event struct {
	Payload   string
	EventType EventType `validate:"oneof='message notification'"`
}

type EventMessage struct {
	Event *Event
	Ids   []uint
}

package websocket


type EventType string

const (
	Message EventType = "message"
	Notification EventType = "notification"
)

type Event struct {
	Payload []byte
	EventType EventType `validate:"oneof='message notification'"`
}


type EventRouter interface {
	Publish(ch <-chan *ChannelMessage)
}
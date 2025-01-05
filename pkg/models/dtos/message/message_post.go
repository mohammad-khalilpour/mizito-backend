package message_dto

import (
	"time"
)

type Message struct {
	Sender uint
	Project uint
	Content []byte
	CreatedAt time.Time
}

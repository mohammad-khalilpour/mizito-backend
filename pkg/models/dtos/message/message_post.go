package message_dto

import (
	"time"
)

type Message struct {
	Project   uint   `json:"project_id"`
	Content   string `json:"content"`
	CreatedAt time.Time
}

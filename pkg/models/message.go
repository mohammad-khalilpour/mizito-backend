package models

import (
	"time"
)

type Message struct {
	Sender    uint
	Project   Project
	Content   []byte
	CreatedAt time.Time
}

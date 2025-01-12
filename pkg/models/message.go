package models

import (
	"time"
)

type Message struct {
	Sender    uint
	ProjectID uint
	Project   Project
	Content   []byte
	CreatedAt time.Time
}

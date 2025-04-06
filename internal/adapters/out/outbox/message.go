package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID
	Name        string
	Payload     []byte
	OccurredAt  time.Time
	ProcessedAt *time.Time
}

func (Message) TableName() string {
	return "outbox"
}

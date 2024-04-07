package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Uuid       uuid.UUID
	Created_at time.Time
	Updated_at time.Time
}

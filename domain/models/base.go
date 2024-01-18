package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModelTable struct {
	Uuid       uuid.UUID
	Created_at time.Time
	Updated_at time.Time
}

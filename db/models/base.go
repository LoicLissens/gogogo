package db

import (
	"time"

	"github.com/google/uuid"
)

type BaseModelTable struct {
	Uuid       uuid.UUID `db:"uuid" sql_properties:"UUID PRIMARY KEY"`
	Created_at time.Time `db:"created_at" sql_properties:"TIMESTAMP NOT NULL"`
	Updated_at time.Time `db:"updated_at" sql_properties:"TIMESTAMP NOT NULL"`
}

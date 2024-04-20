package repositories

import (
	"jiva-guildes/domain/models/guilde"

	"github.com/google/uuid"
)

type GuildeRepository interface {
	GetByUUID(connectionPool interface{}, uuid uuid.UUID, tableName string, schema string) (guilde.Guilde, error)
	Save(connectionPool interface{}, tableName string, schema string, entity interface{}) (guilde.Guilde, error)
}

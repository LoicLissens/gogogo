package repositories

import (
	"jiva-guildes/domain/models/guilde"

	"github.com/google/uuid"
)

type GuildeRepository interface {
	GetByUUID(uuid uuid.UUID) (guilde.Guilde, error)
	Save(guilde.Guilde) (guilde.Guilde, error)
	Delete(uuid uuid.UUID) error
}

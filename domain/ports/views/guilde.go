package views

import (
	"jiva-guildes/domain/ports/views/dtos"

	"github.com/google/uuid"
)

type GuildeView interface {
	Fetch(uuid uuid.UUID) dtos.GuildeViewDTO
	List(page int, limit int)
}

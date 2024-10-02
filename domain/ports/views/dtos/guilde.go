package dtos

import (
	"time"

	"github.com/google/uuid"
)

type GuildeViewDTO struct {
	Uuid          uuid.UUID  `json:"uuid" validate:"required"`
	Name          string     `json:"name" validate:"required"`
	Img_url       string     `json:"imgUrl" validate:"required,url"`
	Page_url      string     `json:"pageUrl" validate:"required,url"`
	Created_at    time.Time  `json:"createdAt" validate:"required"`
	Updated_at    time.Time  `json:"updatedAt" validate:"required"`
	Exists        bool       `json:"exists" validate:"required"`
	Validated     bool       `json:"validated" validate:"required"`
	Active        *bool      `json:"active"`
	Creation_date *time.Time `json:"creationDate" validate:"datetime"`
}
type GuildeListViewDTO struct {
	Items   []GuildeViewDTO
	NbItems int
}

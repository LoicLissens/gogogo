package commands

import (
	"time"

	"github.com/google/uuid"
)

type UpdateGuildeCommand struct {
	Uuid         uuid.UUID `json:"uuid" validate:"required,uuid"`
	Name         string    `json:"name"`
	Img_url      string    `json:"img_url"`
	Page_url     string    `json:"page_url"`
	Exists       *bool     `json:"exists"`
	Active       *bool     `json:"active"`
	CreationDate time.Time `json:"creation_date"`
	Validated    *bool     `json:"validated"`
}

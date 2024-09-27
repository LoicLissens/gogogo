package commands

import "github.com/google/uuid"

type DeleteGuildeCommand struct {
	Uuid uuid.UUID `json:"uuid" validate:"required,uuid"`
}

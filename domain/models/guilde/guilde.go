package models

import (
	"time"

	"github.com/google/uuid"
)

type Guilde struct {
	Uuid       uuid.UUID
	Name       string
	Img_url    string
	Page_url   string
	Created_at time.Time
	Updated_at time.Time
}

func New(name string, img_url string, page_url string) *Guilde {
	guilde := &Guilde{
		Uuid:       uuid.New(),
		Name:       name,
		Img_url:    img_url,
		Page_url:   page_url,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	return guilde
}

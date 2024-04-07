package guilde

import (
	"time"

	"jiva-guildes/domain/models"

	"github.com/google/uuid"
)

type Guilde struct {
	models.BaseModel
	Name     string
	Img_url  string
	Page_url string
}

func New(name string, img_url string, page_url string) *Guilde {
	guilde := &Guilde{
		Name:     name,
		Img_url:  img_url,
		Page_url: page_url,
	}
	guilde.Uuid = uuid.New()
	guilde.Created_at = time.Now()
	guilde.Updated_at = time.Now()

	return guilde
}

package commands

import "time"

type CreateGuildeCommand struct {
	Name          string     `json:"name" validate:"required"`
	Img_url       string     `json:"img_url" validate:"omitempty,url"`
	Page_url      string     `json:"page_url" validate:"omitempty,url"`
	Exists        bool       `json:"exists"`
	Active        *bool      `json:"active"`
	Creation_date *time.Time `json:"creation_date"`
}

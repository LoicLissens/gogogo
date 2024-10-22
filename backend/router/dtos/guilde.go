package dtos

import (
	"jiva-guildes/domain/ports/views"
	viewdtos "jiva-guildes/domain/ports/views/dtos"

	"net/url"
	"strconv"
	"time"
)

type CreateGuildeInput struct {
	Name          string     `json:"name" form:"name" query:"name"`
	Img_url       string     `json:"img_url" form:"img_url" query:"img_url"`
	Page_url      string     `json:"page_url" form:"page_url" query:"page_url"`
	Exists        bool       `json:"exists" form:"exists" query:"exists"`
	Active        *bool      `json:"active" form:"active" query:"active"`
	Creation_date *time.Time `json:"creation_date" form:"creation_date" query:"creation_date"`
}

type ListGuildesInput struct {
	Page           int                  `query:"page"`
	Limit          int                  `query:"limit"`
	OrderingMethod views.OrderingMethod `query:"ordering_method"`

	OrderBy           views.OrderByGuilde `query:"order_by"`
	Name              string              `query:"name"`
	Exists            *bool               `query:"exists"`
	Validated         *bool               `query:"validated"`
	Active            *bool               `query:"active"`
	CreationDateSince time.Time           `query:"creation_date_since"`
	CreationDateUntil time.Time           `query:"creation_date_until"`
}
type ListeGuildePageData struct {
	Lang        string
	Title       string
	Items       []viewdtos.GuildeViewDTO
	NbItems     int
	CurrentPage int
	TotalPages  int
	CurrentURL  string
}

func (d ListeGuildePageData) GetNextPage() string {
	url, err := url.Parse(d.CurrentURL)
	if err != nil {
		panic(err)
	}
	q := url.Query()
	q.Set("page", strconv.Itoa(d.CurrentPage+1))
	url.RawQuery = q.Encode()
	return url.String()
}
func (d ListeGuildePageData) GetPrevPage() string {
	url, err := url.Parse(d.CurrentURL)
	if err != nil {
		panic(err)
	}
	q := url.Query()
	q.Set("page", strconv.Itoa(d.CurrentPage-1))
	url.RawQuery = q.Encode()
	return url.String()
}

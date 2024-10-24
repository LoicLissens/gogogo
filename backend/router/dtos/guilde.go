package dtos

import (
	"jiva-guildes/domain/ports/views"
	viewdtos "jiva-guildes/domain/ports/views/dtos"

	"net/url"
	"strconv"
	"time"
)

type CustomTime struct {
	*time.Time
}

// needed to parse time from query params by default echo use time.RFC3339 (2020-12-09T16:09:53+00:00) but browers use 2020-12-09
// UnmarshalParam is complient with echo.ParamUnmarshaler interface
func (ct *CustomTime) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	date, err := time.Parse("2006-01-02", param)
	if err != nil {
		return err
	}
	ct.Time = &date
	return nil
}

type CreateGuildeInput struct {
	Name          string      `json:"name" form:"name" query:"name"`
	Img_url       string      `json:"img_url" form:"img_url" query:"img_url"`
	Page_url      string      `json:"page_url" form:"page_url" query:"page_url"`
	Exists        bool        `json:"exists" form:"exists" query:"exists"`
	Active        *bool       `json:"active" form:"active" query:"active"`
	Creation_date *CustomTime `json:"creation_date" form:"creation_date" query:"creation_date"`
}

func (c *CreateGuildeInput) GetCreationDate() *time.Time {
	if c.Creation_date != nil {
		return c.Creation_date.Time
	} else {
		return nil
	}
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

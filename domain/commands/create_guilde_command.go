package commands

type CreateGuildeCommand struct {
	Name     string `json:"name" validate:"required"`
	Img_url  string `json:"img_url" validate:"url"`
	Page_url string `json:"page_url" validate:"required,url"`
}

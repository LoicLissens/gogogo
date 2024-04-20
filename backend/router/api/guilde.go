package api

import (
	"github.com/labstack/echo/v4"
)

type GuildeInput struct {
	Name     string `json:"name" validate:"required"`
	Img_url  string `json:"img_url" validate:"datauri"`
	Page_url string `json:"page_url" validate:"required,datauri"`
}

func InitGuildeApiRoutes(e *echo.Echo) {
	api := e.Group("/api")
	api.GET("/guildes", getGuildes)
	api.GET("/guildes/:id", getGuilde)
	api.POST("/guildes", createGuilde)
}

func getGuildes(c echo.Context) error {
	return nil //Todo shoudl return guildes with pagination and tout le tralala
}

func getGuilde(c echo.Context) error {
	return nil
}

func createGuilde(c echo.Context) error {
	return nil
}

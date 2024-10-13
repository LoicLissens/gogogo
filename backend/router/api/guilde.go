package api

import (
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/utils"
	"jiva-guildes/domain/commands"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type CreateGuildeInput struct {
	Name          string     `json:"name" form:"name" query:"name"`
	Img_url       string     `json:"img_url" form:"img_url" query:"img_url"`
	Page_url      string     `json:"page_url" form:"page_url" query:"page_url"`
	Exists        bool       `json:"exists" form:"exists" query:"exists"`
	Active        *bool      `json:"active" form:"active" query:"active"`
	Creation_date *time.Time `json:"creation_date" form:"creation_date" query:"creation_date"`
}

var ServiceManager = backend.ServiceManager

func InitGuildeApiRoutes(e *echo.Echo) {
	api := e.Group("/api")
	api.GET("/guildes", getGuildes)
	api.GET("/guildes/:id", getGuilde)
	api.POST("/guildes", createGuilde)
}

func getGuildes(c echo.Context) error {
	return nil
}

func getGuilde(c echo.Context) error {
	return nil
}

func createGuilde(c echo.Context) error {
	g := new(CreateGuildeInput)
	if err := c.Bind(g); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request body")
	}
	cmd := commands.CreateGuildeCommand{
		Name:          g.Name,
		Img_url:       g.Img_url,
		Page_url:      g.Page_url,
		Exists:        g.Exists,
		Active:        g.Active,
		Creation_date: g.Creation_date,
	}

	if err := backend.Validate.Struct(cmd); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	guilde, err := ServiceManager.CreateGuildeHandler(cmd)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.PostMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.JSON(http.StatusCreated, guilde)
}

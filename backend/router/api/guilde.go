package api

import (
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/dtos"
	"jiva-guildes/backend/router/utils"
	"jiva-guildes/domain/commands"
	"net/http"

	"github.com/labstack/echo/v4"
)

var serviceManager = backend.ServiceManager

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
	g := new(dtos.CreateGuildeInput)
	if err := c.Bind(g); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request body")
	}
	if err := c.Validate(g); err != nil {
		return err
	}
	cmd := commands.CreateGuildeCommand{
		Name:          g.Name,
		Img_url:       g.Img_url,
		Page_url:      g.Page_url,
		Exists:        g.Exists,
		Active:        g.Active,
		Creation_date: g.GetCreationDate(),
	}

	if err := backend.Validate.Struct(cmd); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	} //TODO: should go in the service layer
	guilde, err := serviceManager.CreateGuildeHandler(cmd)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.PostMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.JSON(http.StatusCreated, guilde)
}

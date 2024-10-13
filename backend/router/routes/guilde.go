package routes

import (
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var serviceManager = backend.ServiceManager
var viewsManager = backend.ViewsManager

func InitGuildeRoutes(e *echo.Echo) {
	api := e.Group("")
	api.GET("/guildes/:uuid", getGuilde)
}

func getGuilde(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request body")
	}
	g, err := viewsManager.Guilde().Fetch(uuid)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.GetMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.JSON(http.StatusOK, g)
}

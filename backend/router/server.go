package router

import (
	"net/http"

	"jiva-guildes/backend/router/api"

	"github.com/labstack/echo/v4"
)

func Serve() { //TODO call teardown in case of crash or shutdown
	e := echo.New()
	api.InitGuildeApiRoutes(e)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

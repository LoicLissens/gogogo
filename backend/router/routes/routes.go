package routes

import "github.com/labstack/echo/v4"

func IsHTMXRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

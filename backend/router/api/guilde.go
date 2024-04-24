package api

import (
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/utils"
	"jiva-guildes/domain/commands"

	"github.com/labstack/echo/v4"
)

type GuildeInput struct {
	Name     string `json:"name" form:"name" query:"name"`
	Img_url  string `json:"img_url" form:"img_url" query:"img_url"`
	Page_url string `json:"page_url" form:"page_url" query:"page_url"`
}

// HandleInput is a function that binds the request body to the data struct, validates it
// and instanciate a command. Should always be used for POST requests.
func HandleInput(c echo.Context, data interface{}) error {
	if err := c.Bind(data); err != nil { //TODO Maybe create an utility function to bind/validate(/return a cmd) to use in every request
		return echo.NewHTTPError(utils.StatusBadRequest, "Failed to parse request body")
	}
	return nil
}

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
	g := new(GuildeInput)
	if err := c.Bind(g); err != nil {
		return echo.NewHTTPError(utils.StatusBadRequest, "Failed to parse request body")
	}
	cmd := commands.CreateGuildeCommand{
		Name:     g.Name,
		Img_url:  g.Img_url,
		Page_url: g.Page_url,
	}

	if err := backend.Validate.Struct(cmd); err != nil {
		return echo.NewHTTPError(utils.StatusUnprocessable, err.Error())
	}
	guilde, err := backend.ServiceManager.CreateGuildeHandler(cmd)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.PostMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.JSON(utils.StatusCreated, guilde)
}

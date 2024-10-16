package routes

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/utils"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views"
	viewdtos "jiva-guildes/domain/ports/views/dtos"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var viewsManager = backend.ViewsManager
var serviceManager = backend.ServiceManager

func InitGuildeRoutes(e *echo.Echo) {
	api := e.Group("")
	api.GET("/guildes/:uuid", getGuilde)
	api.GET("/guildes", listGuildes)
	api.DELETE("/guildes/:uuid", deleteGuilde)
}

func getGuilde(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse path parameter")
	}
	g, err := viewsManager.Guilde().Fetch(uuid)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.GetMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.JSON(http.StatusOK, g)
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
	CreationDateSince time.Time           `query:"creation_date_since" validate:"datetime"`
	CreationDateUntil time.Time           `query:"creation_date_until" validate:"datetime"`
}
type ListeGuildePageData struct {
	Lang        string
	Title       string
	Items       []viewdtos.GuildeViewDTO
	NbItems     int
	CurrentPage int
	TotalPages  int
}

func (d ListeGuildePageData) GetNextPage() int {
	return d.CurrentPage + 1
}
func (d ListeGuildePageData) GetPrevPage() int {
	return d.CurrentPage - 1
}

func listGuildes(c echo.Context) error {
	input := new(ListGuildesInput)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse query parameters")
	}
	// if err := c.Validate(input); err != nil { //TODO: validate the datetime make the page crash
	// 	return err
	// }
	var limit int // EUurk disgusting !!!
	if input.Limit == 0 {
		limit = 10
	} else {
		limit = input.Limit
	}

	guildes, err := viewsManager.Guilde().List(views.ListGuildesViewOpts{
		Page:              input.Page,
		Limit:             limit,
		OrderingMethod:    input.OrderingMethod,
		OrderBy:           input.OrderBy,
		Name:              input.Name,
		Exists:            input.Exists,
		Validated:         input.Validated,
		Active:            input.Active,
		CreationDateSince: input.CreationDateSince,
		CreationDateUntil: input.CreationDateUntil,
	})
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.GetMethod)
		c.Logger().Error(fmt.Sprintf("Error while fetching guildes: %s", message))
		return echo.NewHTTPError(code, message)
	}
	data := ListeGuildePageData{
		Lang:        "fr",
		Title:       "Liste",
		Items:       guildes.Items,
		NbItems:     guildes.NbItems,
		CurrentPage: input.Page,
		TotalPages:  guildes.NbItems / limit,
	}

	var templateName string
	if utils.IsHTMXRequest(c) {
		templateName = "display"
	} else {
		templateName = "guildes.html"
	}
	return c.Render(http.StatusOK, templateName, data)
}
func deleteGuilde(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse path parameter")
	}
	err = serviceManager.DeleteGuildeHandler(commands.DeleteGuildeCommand{Uuid: uuid})
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.DeleteMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.NoContent(http.StatusOK)
}

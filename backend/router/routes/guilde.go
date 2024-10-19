package routes

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/utils"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views"
	viewdtos "jiva-guildes/domain/ports/views/dtos"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var viewsManager = backend.ViewsManager
var serviceManager = backend.ServiceManager

func InitGuildeRoutes(e *echo.Echo) {
	api := e.Group("")
	api.GET("/guildes", listGuildes)
	api.GET("/guildes/edit/:uuid", editGuilde)
	api.GET("/guildes/row/:uuid", getRowGuilde)
	api.DELETE("/guildes/:uuid", deleteGuilde)
	api.PATCH("/guildes/:uuid", patchGuilde)
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

// ROUTES HANDLERS

func listGuildes(c echo.Context) error {
	input := new(ListGuildesInput)
	if err := c.Bind(input); err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing query parameters: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse query parameters")
	}
	// if err := c.Validate(input); err != nil { //TODO: validate the datetime make the page crash
	// 	return err
	// }
	page, limit := utils.GetPageAndLimit(input.Page, input.Limit)

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
		CurrentPage: page,
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

func editGuilde(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse path parameter")
	}
	g, err := viewsManager.Guilde().Fetch(uuid)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.GetMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.Render(http.StatusOK, "edit-guilde-row", g)
}
func getRowGuilde(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse path parameter")
	}
	g, err := viewsManager.Guilde().Fetch(uuid)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.GetMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.Render(http.StatusOK, "guilde-row", g)
}
func patchGuilde(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("uuid"))
	name := c.FormValue("name")
	validated := c.FormValue("validated")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse path parameter")
	}
	var command commands.UpdateGuildeCommand
	command.Name = name
	command.Uuid = uuid
	validatedBool, err := strconv.ParseBool(validated)
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing form value: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse form value")
	}
	command.Validated = &validatedBool
	fmt.Println("validatedBool", validatedBool)
	updatedG, err := serviceManager.UpdateGuildeHandler(command)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.PutMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.Render(http.StatusOK, "guilde-row", updatedG.ToViewDTO())
}

package routes

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/backend/router/dtos"
	"jiva-guildes/backend/router/utils"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views"
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
	api.POST("/guildes", createGuilde)
}

// ROUTES HANDLERS

func listGuildes(c echo.Context) error {
	input := new(dtos.ListGuildesInput)
	if err := c.Bind(input); err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing query parameters: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse query parameters")
	}
	if err := c.Validate(input); err != nil {
		return err
	}
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
	data := dtos.ListeGuildePageData{
		Lang:        "fr",
		Title:       "Liste",
		Items:       guildes.Items,
		NbItems:     guildes.NbItems,
		CurrentPage: page,
		TotalPages:  guildes.NbItems / limit,
		CurrentURL:  c.Request().URL.String(),
	}

	var templateName string
	if utils.IsHTMXRequest(c) {
		templateName = "display-guildes"
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
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse path parameter")
	}
	name := c.FormValue("name")
	creationDate := c.FormValue("creation-date")
	validated := c.FormValue("validated")
	exists := c.FormValue("exists")
	active := c.FormValue("active")

	var parsedCreationDate time.Time
	if creationDate != "" {
		parsedCreationDate, err = time.Parse("2006-01-02", creationDate)
		if err != nil {
			c.Logger().Error(fmt.Sprintf("Error while parsing form value: %s", err))
			return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse form value")
		}
	} else {
		parsedCreationDate = time.Time{}
	}
	parsedActive, err := strconv.ParseBool(active)
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing form value: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse form value")
	}
	parsedExists, err := strconv.ParseBool(exists)
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing form value: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse form value")
	}
	validatedBool, err := strconv.ParseBool(validated)
	if err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing form value: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse form value")
	}

	var command commands.UpdateGuildeCommand
	command.Name = name
	command.Uuid = uuid
	command.CreationDate = parsedCreationDate
	command.Validated = &validatedBool
	command.Exists = &parsedExists
	command.Active = &parsedActive

	updatedG, err := serviceManager.UpdateGuildeHandler(command)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.PutMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.Render(http.StatusOK, "guilde-row", updatedG.ToViewDTO())
}

func createGuilde(c echo.Context) error {
	g := new(dtos.CreateGuildeInput)
	if err := c.Bind(g); err != nil {
		c.Logger().Error(fmt.Sprintf("Error while parsing request body: %s", err))
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request body")
	}
	img, err := c.FormFile("image")
	if err != nil {
		if err.Error() != "http: no such file" {
			c.Logger().Error(fmt.Sprintf("Error while parsing img file: %s", err))
			code, message := utils.ErrorCodeMapper(err, utils.PostMethod)
			return echo.NewHTTPError(code, message)
		}
	}
	if img != nil {
		fmt.Println("img", img.Size)
	}
	if err := c.Validate(g); err != nil {
		c.Logger().Error(fmt.Sprintf("Error while validating request body: %s", err))
		code, message := utils.ErrorCodeMapper(err, utils.PostMethod)
		return echo.NewHTTPError(code, message)
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
	}
	guilde, err := serviceManager.CreateGuildeHandler(cmd)
	if err != nil {
		code, message := utils.ErrorCodeMapper(err, utils.PostMethod)
		return echo.NewHTTPError(code, message)
	}
	return c.Render(http.StatusCreated, "creation-success", guilde.ToViewDTO())
}

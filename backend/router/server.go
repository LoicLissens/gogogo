package router

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"time"

	"jiva-guildes/backend"
	"jiva-guildes/backend/router/api"
	"jiva-guildes/backend/router/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed views/*.html
var tmplFS embed.FS

//go:embed all:static
var assetsFS embed.FS

type (
	CustomValidator struct {
		validator *validator.Validate
	}
	Template struct { // TODO use template cache ? https://www.youtube.com/watch?v=JbtHT1-vAfA
		tmpl *template.Template
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func newTemplate() *Template {
	funcMap := template.FuncMap{
		"GetDate": func(date *time.Time) string {
			if date == nil {
				return ""
			}
			return date.Format("01-02-2006")
		},
	}
	// Using Must to ensure that the template is loaded at the start of the application as it will panic if template is not valid
	// same as doing this: template.New(tmplFile).ParseFiles(tmplFile) then checking the error returned and panic if so
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFS(tmplFS, "views/*.html"))
	return &Template{
		tmpl: tmpl,
		// parse all html files in the views folder in a *Template collection
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
	// ExecuteTemplate is used to apply a specific named template in the collection to the data
}
func Serve() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: backend.Validate}
	e.Use(middleware.Logger())

	subFS := echo.MustSubFS(assetsFS, "static")
	e.StaticFS("/static/*", subFS)

	e.Renderer = newTemplate()
	api.InitGuildeApiRoutes(e)
	routes.InitGuildeRoutes(e)
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{"lang": "en", "title": "Index gogogo"})
	})
	e.Logger.Fatal(e.Start(":1323"))
}

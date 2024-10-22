package utils

import (
	"errors"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/settings"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	JsonType          = "application/json"
	ZipType           = "application/zip"
	AppXMLType        = "application/xml"
	XmlType           = "text/xml"
	SldType           = "application/vnd.ogc.sld+xml"
	ContentTypeHeader = "Content-Type"
	AcceptHeader      = "Accept"
	GetMethod         = "GET"
	PutMethod         = "PUT"
	PostMethod        = "POST"
	DeleteMethod      = "DELETE"
)

func ErrorCodeMapper(err error, method string) (int, string) {
	switch err.(type) {
	case nil:
		if method == PostMethod {
			return http.StatusCreated, ""
		}
		return http.StatusOK, ""
	case customerrors.ErrorNotFound:
		return http.StatusNotFound, err.Error()
	case customerrors.ErrorAlreadyExists:
		return http.StatusUnprocessableEntity, err.Error()
	default:
		return http.StatusInternalServerError, errors.New("internal server error").Error()
	}
}
func IsHTMXRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}
func GetPageAndLimit(page int, limite int) (int, int) {
	if page == 0 {
		page = settings.AppSettings.DEFAULT_PAGE
	}
	if limite == 0 {
		limite = settings.AppSettings.DEFAULT_PAGE_LIMIT
	}
	return page, limite
}
func GetDateForTemplate(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format("01-02-2006")
}

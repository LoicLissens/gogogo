package utils

import (
	customerrors "jiva-guildes/domain/custom_errors"
)

const (
	StatusOk            = 200
	StatusCreated       = 201
	StatusNotAllowed    = 405
	StatusForbidden     = 403
	StatusUnprocessable = 422
	StatusInternalError = 500
	StatusNotFound      = 404
	StatusUnauthorized  = 401
	StatusBadRequest    = 400
	JsonType            = "application/json"
	ZipType             = "application/zip"
	AppXMLType          = "application/xml"
	XmlType             = "text/xml"
	SldType             = "application/vnd.ogc.sld+xml"
	ContentTypeHeader   = "Content-Type"
	AcceptHeader        = "Accept"
	GetMethod           = "GET"
	PutMethod           = "PUT"
	PostMethod          = "POST"
	DeleteMethod        = "DELETE"
)

func ErrorCodeMapper(err error, method string) (int, string) {
	switch err.(type) {
	case nil:
		if method == PostMethod {
			return StatusCreated, ""
		}
		return StatusOk, ""
	case customerrors.ErrorNotFound:
		return StatusNotFound, err.Error()
	case customerrors.ErrorAlreadyExists:
		return StatusUnprocessable, err.Error()
	default:
		return StatusInternalError, err.Error()
	}
}

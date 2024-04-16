Use a PostgresQL docker containet for the DB:

docker run -d --name jiva-g -p 127.0.0.1:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=jiva-g postgres

go install github.com/cosmtrek/air@latest => live reload

go install golang.org/x/tools/gopls@latest

[text](https://templ.guide/) => idea of templating

if issue with importing external package, run : go mod tidy

Sanetisation et honeypot ?

use this const (
	statusOk            = 200
	statusCreated       = 201
	statusNotAllowed    = 405
	statusForbidden     = 403
	statusInternalError = 500
	statusNotFound      = 404
	statusUnauthorized  = 401
	jsonType            = "application/json"
	zipType             = "application/zip"
	appXMLType          = "application/xml"
	xmlType             = "text/xml"
	sldType             = "application/vnd.ogc.sld+xml"
	contentTypeHeader   = "Content-Type"
	acceptHeader        = "Accept"
	getMethod           = "GET"
	putMethod           = "PUT"
	postMethod          = "POST"
	deleteMethod        = "DELETE"
)
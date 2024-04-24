Run with hot reloading wit *air*:
- install globaly air: `go install github.com/cosmtrek/air@latest`
- set an alias for air: `alias air='$(go env GOPATH)/bin/air'`
- if not implemented: `air init`
- launch with: `air`


Use a PostgresQL docker containet for the DB:

```docker run -d --name container-name -p 127.0.0.1:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=dn-name postgres```

If issue with importing external package, run : `go mod tidy`


Run with hot reloading wit *air*:
- install globaly air: `go install github.com/cosmtrek/air@latest`
- set an alias for air: `alias air='$(go env GOPATH)/bin/air'`
- if not implemented: `air init`
- launch with: `air`

Run with cli :
`go run . -cli`

Run tests:
cd in the directory `go test .` eg fromthe adapters : ` go test ./db/repositories`

Use a PostgresQL docker containet for the DB:

```docker run -d --name container-name -p 127.0.0.1:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=dn-name postgres```

If issue with importing external package, run : `go mod tidy`

To make a copy of the DB run :
`docker exec -i container-name /bin/bash -c "PGPASSWORD=pass pg_dump --username username dbname" > dump.sql`
Then to reaply a copy to the db :
`docker exec -i container-name /bin/bash -c "PGPASSWORD=pass psql --username username dbname" < dump.sql`
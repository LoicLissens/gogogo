Use a PostgresQL docker containet for the DB:

docker run -d --name jiva-g -p 127.0.0.1:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=jiva-g postgres

go install github.com/cosmtrek/air@latest => live reload

go install golang.org/x/tools/gopls@latest
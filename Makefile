include .env

init-db:
	docker run -d --name $(DATABASE_CONTAINER_NAME) -p 127.0.0.1\:5432:5432 -e POSTGRES_USER=$(DATABASE_USER) -e POSTGRES_PASSWORD=$(DATABASE_PASSWORD) -e POSTGRES_DB=$(DATABASE_NAME) postgres

rm-db:
	docker stop $(DATABASE_CONTAINER_NAME) && docker rm $(DATABASE_CONTAINER_NAME)

reset-db: rm-db init-db

tests:
	go test ./... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''

DATE = $(shell date +%Y-%m-%d)
SRC_FILE = 'dump.sql'
save-db:
	docker exec -i $(DATABASE_CONTAINER_NAME) /bin/bash -c "PGPASSWORD=$(DATABASE_PASSWORD) pg_dump --username $(DATABASE_USER) $(DATABASE_NAME)" >  $(DB_DUMP_FOLDER)/$(SRC_FILE)
	cp dump.sql $(DB_DUMP_FOLDER)/$(basename $(SRC_FILE))_$(DATE)$(suffix $(SRC_FILE))

restore-db:
	docker exec -i $(DATABASE_CONTAINER_NAME) /bin/bash -c "PGPASSWORD=$(DATABASE_PASSWORD) psql --username $(DATABASE_USER) $(DATABASE_NAME)" <  $(DB_DUMP_FOLDER)/$(SRC_FILE)

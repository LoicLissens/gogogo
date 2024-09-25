init-db:
	docker run -d --name jiva-g -p 127.0.0.1\:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=jiva-g postgres

rm-db:
	docker stop jiva-g && docker rm jiva-g

reset-db: rm-db init-db

tests:
	go test ./... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''

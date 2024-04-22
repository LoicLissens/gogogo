Use a PostgresQL docker containet for the DB:

docker run -d --name jiva-g -p 127.0.0.1:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=jiva-g postgres

Hotreloading can be implemented with: `air init`
If issue with importing external package, run : go mod tidy

Try to test with -race = go test -race


TODO:

go install golang.org/x/tools/gopls@latest

[text](https://templ.guide/) => idea of templating


Sanetisation et honeypot ?

use cancelations the net  ninja tuto ? => gracefull shutdown

use https://github.com/mdempsky/maligned?ref=hackernoon.com (maligned) for perf issue ?


organise project with cmd and internal

https://www.reddit.com/r/golang/comments/1ad14so/i_used_reflection_and_i_feel_dirty/ => middleware in routing
https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx => improve row scanning

https://github.com/go-playground/assert => to use for testing

Add echo lock ?
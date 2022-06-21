.PHONY: build clean test sqlc db/run db/migrate db/kill  

build:
	go build

clean:
	-rm sqlc

test:
	go test -v -cover ./...

sqlc:
	sqlc generate

db/run:
	docker run \
		--name sqlc-mysql \
		-e MYSQL_ROOT_PASSWORD=example \
		-e MYSQL_DATABASE=sqlc \
		-e MYSQL_ROOT_HOST=% \
		-p 3306:3306 \
		--detach --rm \
		mysql/mysql-server:latest \
		--character-set-server=utf8mb4

db/migrate:
	migrate -path ./db/migrations/ -database 'mysql://root:example@tcp(127.0.0.1:3306)/sqlc' -verbose up

db/kill:
	docker kill sqlc-mysql


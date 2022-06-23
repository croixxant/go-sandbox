COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

build:
	go build

test:
	go test -v -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

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

.PHONY: build test cover sqlc db/run db/migrate db/kill
.PHONY: build clean server test sqlc mock db/run db/migrate db/kill  

build:
	go build

clean:
	-rm gomock

server:
	go run main.go

test:
	go test -v -cover ./...

sqlc:
	sqlc generate

mock:
	mockgen -destination db/mock/store.go github.com/croixxant/go-sandbox/gomock/db Store

db/run:
	docker run \
		--name gomock-mysql \
		-e MYSQL_ROOT_PASSWORD=example \
		-e MYSQL_DATABASE=gomock \
		-e MYSQL_ROOT_HOST=% \
		-p 3306:3306 \
		--detach --rm\
		mysql/mysql-server:latest \
		--character-set-server=utf8mb4

db/migrate:
	migrate -path ./db/migrations/ -database 'mysql://root:example@tcp(127.0.0.1:3306)/gomock' -verbose up

db/kill:
	docker kill gomock-mysql
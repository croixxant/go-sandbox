.PHONY: build clean test sqlc mock db/run db/migrate db/kill  

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
	mockgen -destination db/mock/store.go github.com/croixxant/golang-examples/gomock/db Store

docker/init: docker/buildapp docker/createnetwork

docker/run: docker/rundb docker/runapp

docker/buildapp:
	docker build -t gomock:latest .

docker/createnetwork:
	docker network create gomock-network

docker/rundb:
	docker run \
		--name gomock-mysql \
		-e MYSQL_ROOT_PASSWORD=example \
		-e MYSQL_DATABASE=gomock \
		-p 3306:3306 \
		--detach --rm \
		--network gomock-network \
		mysql:8.0

docker/runapp:
	docker run \
		--name gomock \
		-p 8080:8080 \
		-e DB_SOURCE="root:example@tcp(gomock-mysql:3306)/gomock" \
		--network gomock-network \
		--detach --rm \
		--entrypoint="/app/wait-for.sh" \
		gomock:latest gomock-mysql:3306 --timeout=60 -- "/app/start.sh" "/app/main"

docker/migrate:
	migrate -path ./db/migrations/ -database 'mysql://root:example@tcp(127.0.0.1:3306)/gomock' -verbose up

docker/clean: docker/kill docker/rmnetwork docker/rmimage

docker/kill:
	docker kill gomock-mysql gomock

docker/rmnetwork:
	docker network rm gomock-network

docker/rmimage:
	docker image rm gomock:latest
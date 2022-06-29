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

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go github.com/croixxant/go-sandbox/simplebank/db Store

proto:
	rm -f proto/*.pb.go \
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=./doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simplebank \
    	proto/*.proto

# MYSQL_ROOT_HOST: allow access from local machine
# --sql-mode: ignore STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO
db/run:
	docker run \
		--name simplebank-mysql \
		--network simplebank-network \
		-e MYSQL_ROOT_PASSWORD=example \
		-e MYSQL_DATABASE=simplebank \
		-e MYSQL_ROOT_HOST=% \
		-p 3306:3306 \
		--detach --rm \
		mysql/mysql-server:latest \
		--character-set-server=utf8mb4 \
		--sql-mode=ONLY_FULL_GROUP_BY,NO_ENGINE_SUBSTITUTION

db/migrateup:
	migrate -path ./db/migrations/ -database 'mysql://root:example@tcp(127.0.0.1:3306)/simplebank' -verbose up 1

db/migratedown:
	migrate -path ./db/migrations/ -database 'mysql://root:example@tcp(127.0.0.1:3306)/simplebank' -verbose down 1

db/kill:
	docker kill simplebank-network

network:
	docker network create simplebank-network


.PHONY: build test cover sqlc server mock proto db/run db/migrateup db/migratedown db/kill docker/network
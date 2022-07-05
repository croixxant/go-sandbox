COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
DB_SOURCE=root:example@tcp(127.0.0.1:3306)
APP_NAME=sandbox
PACKAGE_NAME=github.com/croixxant/go-sandbox

.PHONY: gin grpc test cover sqlc mock proto migrateup migrateupall migratedown migratedownall rundb killdb

gin:
	go build -o bin/ ./cmd/gin/main.go

grpc:
	go build -o bin/ ./cmd/grpc/main.go

test:
	go test -v -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

sqlc:
	sqlc generate -f ./sqlc.yaml

mock:
	mockgen -destination repo/mock/repository.go $(PACKAGE_NAME)/usecase/repo Repository

proto:
	rm -f controller/grpc/internal/*.pb.go controller/grpc/internal/*.pb.gw.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=controller/grpc/proto --go_out=controller/grpc/internal --go_opt=paths=source_relative \
    	--go-grpc_out=controller/grpc/internal --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=controller/grpc/internal --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=./doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=$(APP_NAME) \
    	controller/grpc/proto/*.proto

migrateup:
	migrate -path ./migration/ -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose up 1

migrateupall:
	migrate -path ./migration/ -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose up

migratedown:
	migrate -path ./migration/ -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose down 1

migratedownall:
	migrate -path ./migration/ -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose down

rundb:
	docker run \
		--name $(APP_NAME)-mysql \
		-e MYSQL_ROOT_PASSWORD=example \
		-e MYSQL_DATABASE=$(APP_NAME) \
		-e MYSQL_ROOT_HOST=% \
		-p 3306:3306 \
		--detach --rm \
		mysql/mysql-server:latest \
		--character-set-server=utf8mb4

killdb:
	docker kill $(APP_NAME)-mysql


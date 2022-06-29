APPNAME=sqlc
DB_DSN=root:example@tcp(127.0.0.1:3306)

.PHONY: sqlc/sqlc sqlc/migrateup sqlc/migratedown sqlc/rundb sqlc/killdb

sqlc/sqlc:
	sqlc generate -f ./configs/sqlc/sqlc.yaml

sqlc/migrateup:
	migrate -path ./pkg/sqlc/migrations/ -database 'mysql://$(DB_DSN)/sqlc' -verbose up

sqlc/migratedown:
	migrate -path ./pkg/sqlc/migrations/ -database 'mysql://$(DB_DSN)/sqlc' -verbose down

sqlc/rundb:
	docker run \
		--name sqlc-mysql \
		-e MYSQL_ROOT_PASSWORD=example \
		-e MYSQL_DATABASE=sqlc \
		-e MYSQL_ROOT_HOST=% \
		-p 3306:3306 \
		--detach --rm \
		mysql/mysql-server:latest \
		--character-set-server=utf8mb4

sqlc/killdb:
	docker kill sqlc-mysql

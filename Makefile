make-test:
	go test ./... -v -cover

DB_NAME = auth
DB_HOST = localhost
DB_PORT = 5432
SSL_MODE = disable

force_db:
	migrate -database postgres://auth:auth@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path db/migrations force

version_db:
	migrate -database postgres://auth:auth@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path db/migrations version

migrate_up:
	migrate -database postgres://auth:auth@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path db/migrations up

migrate_down:
	migrate -database postgres://auth:auth@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path db/migrations down



make-proto:
	cd internal/grpc && \
	rm -f pb/*.pb.go && \
	protoc \
    --proto_path=protos \
    --go_out=pb \
    --go_opt=paths=source_relative \
    --go-grpc_out=pb \
    --go-grpc_opt=paths=source_relative \
    protos/*.proto

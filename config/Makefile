generate-go:
	@protoc \
		--go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
		--go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		config.proto

example-go:
	@go run example/main.go | jq .

defaults-go:
	@go run defaults/main.go | jq .
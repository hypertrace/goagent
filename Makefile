.DEFAULT_GOAL := test

.PHONY: test
test: test-unit test-docker

.PHONY: test-unit
test-unit:
	@go test -count=1 -v -race -cover ./...

.PHONY: docker-test
test-docker:
	@./tests/docker/test.sh

.PHONY: bench
bench:
	go test -v -run - -bench . -benchmem ./...

.PHONY: lint
lint:
	@echo "Running linters..."
	@golangci-lint run ./... && echo "Done."

.PHONY: deps
deps:
	@go get -v -t -d ./...

.PHONY: ci-deps
deps-ci:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

check-examples:
	go build -o ./examples/http_client instrumentation/hypertrace/net/hyperhttp/examples/client/main.go && rm ./examples/http_client
	go build -o ./examples/http_server instrumentation/hypertrace/net/hyperhttp/examples/server/main.go && rm ./examples/http_server
	go build -o ./examples/grpc_client instrumentation/hypertrace/google.golang.org/hypergrpc/examples/client/main.go && rm ./examples/grpc_client
	go build -o ./examples/grpc_server instrumentation/hypertrace/google.golang.org/hypergrpc/examples/server/main.go && rm ./examples/grpc_server
	cd instrumentation/hypertrace/database/hypersql/examples/query; go build -o example main.go && rm ./example
	go build -o ./examples/http_client instrumentation/opentelemetry/net/hyperhttp/examples/client/main.go && rm ./examples/http_client
	go build -o ./examples/http_server instrumentation/opentelemetry/net/hyperhttp/examples/server/main.go && rm ./examples/http_server
	go build -o ./examples/grpc_client instrumentation/opentelemetry/google.golang.org/hypergrpc/examples/client/main.go && rm ./examples/grpc_client
	go build -o ./examples/grpc_server instrumentation/opentelemetry/google.golang.org/hypergrpc/examples/server/main.go && rm ./examples/grpc_server
	go build -o ./examples/http_client instrumentation/opencensus/net/hyperhttp/examples/client/main.go && rm ./examples/http_client
	go build -o ./examples/http_server instrumentation/opencensus/net/hyperhttp/examples/server/main.go && rm ./examples/http_server
	go build -o ./examples/grpc_client instrumentation/opencensus/google.golang.org/hypergrpc/examples/client/main.go && rm ./examples/grpc_client
	go build -o ./examples/grpc_server instrumentation/opencensus/google.golang.org/hypergrpc/examples/server/main.go && rm ./examples/grpc_server

generate-config: # generates config object for Go
	@echo "Compiling the proto file"
	@# use protoc v3.13 and protoc-gen-go v1.25.0
	@cd config/agent-config; protoc --go_out=paths=source_relative:.. config.proto
	@echo "Generating the loaders"
	@cd config; go run cmd/generator/main.go agent-config/config.proto
	@echo "Done."

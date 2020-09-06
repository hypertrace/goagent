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

run-http-client-example:
	go run examples/http/client/main.go

run-http-server-example:
	go run examples/http/server/main.go

run-grpc-client-example:
	go run examples/grpc/client/main.go

run-grpc-server-example:
	go run examples/grpc/server/main.go

check-examples:
	go build -o ./examples/http_client examples/http/client/main.go && rm ./examples/http_client
	go build -o ./examples/http_server examples/http/server/main.go && rm ./examples/http_server
	go build -o ./examples/grpc_client examples/grpc/client/main.go && rm ./examples/grpc_client
	go build -o ./examples/grpc_server examples/grpc/server/main.go && rm ./examples/grpc_server

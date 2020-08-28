.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -race -cover ./...

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
ci-deps:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

run-http-client-example:
	go run examples/http/client/main.go

run-http-server-example:
	go run examples/http/server/main.go

run-grpc-client-example:
	go run examples/grpc/client/main.go

run-grpc-server-example:
	go run examples/grpc/server/main.go

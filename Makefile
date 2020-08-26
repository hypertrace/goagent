test:
	go test -v ./...

bench:
	go test -run=Bench -bench=. ./...

deps:
	go get -v -t -d ./...

run-http-server-example:
	go run example/http/server/main.go

run-grpc-client-example:
	go run examples/grpc/client/main.go

run-grpc-server-example:
	go run examples/grpc/server/main.go

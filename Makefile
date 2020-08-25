test:
	go test -v ./...

bench:
	go test -run=Bench -bench=. ./...

deps:
	go get -v -t -d ./...

run-http-server-example:
	go run example/http/server/main.go

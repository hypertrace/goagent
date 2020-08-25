test:
	go test -v ./...

deps:
	go get -v -t -d ./...

run-http-server-example:
	go run example/http/server/main.go

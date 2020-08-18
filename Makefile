test:
	go test -v ./...

deps:
	go get -v -t -d ./...

run-example:
	go run example/main.go

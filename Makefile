test:
	go test -v ./...

deps:
	GO111MODULE=on go get -v -t -d ./...

run-example:
	go run example/main.go
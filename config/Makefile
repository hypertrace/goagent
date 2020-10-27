generate-go: # generates proto object for Go
	#@protoc --go_out=plugins=grpc,paths=source_relative:./transport config.proto
	@go run cmd/generator/main.go config.proto
	@gofmt -s -w config.ps.go

example-go: # runs the Go example
	@echo "Parsed values:"
	@echo ""
	@go run example/main.go | jq .

defaults-go: # displays the Go defaults
	@go run defaults/main.go | jq .
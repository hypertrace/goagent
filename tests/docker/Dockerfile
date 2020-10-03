FROM golang:1.14

WORKDIR $GOPATH/src/github.com/traceableai/goagent/docker

COPY tests/docker/main.go main.go
COPY sdk/internal/container.go internal/container.go

RUN go get -d -v ./...

RUN go build -o /test main.go

ENTRYPOINT [ "/test" ] 

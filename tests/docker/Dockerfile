FROM golang:1.14

WORKDIR $GOPATH/src/github.com/hypertrace/goagent/docker

COPY tests/docker/main.go main.go
COPY sdk/internal/container/container.go internal/container/container.go

RUN go get -d -v ./...

RUN go build -o /test main.go

ENTRYPOINT [ "/test" ] 

module github.com/hypertrace/goagent

go 1.15

require (
	cloud.google.com/go v0.81.0 // indirect
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/hypertrace/agent-config/gen/go v0.0.0-20220627203712-00a3c5053ec6
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/openzipkin/zipkin-go v0.4.0
	github.com/stretchr/testify v1.7.5
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.31.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.31.0
	go.opentelemetry.io/contrib/propagators/b3 v1.6.0
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.3
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.6.3
	go.opentelemetry.io/otel/exporters/zipkin v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

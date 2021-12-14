module github.com/hypertrace/goagent

go 1.15

require (
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/hypertrace/agent-config/gen/go v0.0.0-20211201093046-a7ed863ff59e
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.25.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.25.0
	go.opentelemetry.io/contrib/propagators/b3 v1.0.0
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.1
	go.opentelemetry.io/otel/exporters/zipkin v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	golang.org/x/net v0.0.0-20211008194852-3b03d305991f // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211008145708-270636b82663 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

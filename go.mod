module github.com/hypertrace/goagent

go 1.14

require (
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.6.2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.22.4
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.16.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.16.0
	go.opentelemetry.io/contrib/propagators v0.16.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200803210538-64077c9b5642 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gotest.tools v2.2.0+incompatible
)

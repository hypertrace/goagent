module github.com/hypertrace/goagent

go 1.14

require (
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	contrib.go.opencensus.io/integrations/ocsql v0.1.6
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.4.2
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/openzipkin/zipkin-go v0.2.3
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.4
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc v0.11.0
	go.opentelemetry.io/contrib/instrumentation/net/http v0.11.0
	go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/exporters/stdout v0.11.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.11.0
	go.opentelemetry.io/otel/sdk v0.11.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200803210538-64077c9b5642 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
)

package opentelemetry // import "github.com/hypertrace/goagent/instrumentation/opentelemetry"

import (
	"net/http"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// Server HTTP metrics.
const (
	meterName = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	// Pseudo of go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#RequestCount since a metric is not
	// created for that one for some reason.(annotated with hypertrace to avoid a duplicate if otel go ever implement
	// their own)
	requestCountCounterName = "hypertrace.http.server.request_count" // Incoming request count total
)

type HttpOperationMetricsHandler struct {
	operationNameGetter func(*http.Request) string
	requestCountCounter instrument.Int64Counter
}

var _ sdk.HttpOperationMetricsHandler = (*HttpOperationMetricsHandler)(nil)

func NewHttpOperationMetricsHandler(nameGetter func(*http.Request) string) sdk.HttpOperationMetricsHandler {
	mp := global.MeterProvider()
	meter := mp.Meter(meterName, metric.WithInstrumentationVersion(otelhttp.SemVersion()))

	// Set up net http metrics
	// RequestCount Counter
	requestCountCounter, err := meter.Int64Counter(requestCountCounterName)
	if err != nil {
		otel.Handle(err)
	}

	return &HttpOperationMetricsHandler{
		operationNameGetter: nameGetter,
		requestCountCounter: requestCountCounter,
	}
}

func (mh *HttpOperationMetricsHandler) AddToRequestCount(n int64, r *http.Request) {
	// Add metrics using the same logic in go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#handler.go
	ctx := r.Context()
	labeler, _ := otelhttp.LabelerFromContext(ctx)
	operationName := mh.operationNameGetter(r)
	attributes := append(labeler.Get(), semconv.HTTPServerMetricAttributesFromHTTPRequest(operationName, r)...)
	mh.requestCountCounter.Add(ctx, n, attributes...)
}

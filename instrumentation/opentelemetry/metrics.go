package opentelemetry // import "github.com/hypertrace/goagent/instrumentation/opentelemetry"

import (
	// "bytes"
	// "encoding/base64"
	// "io"
	// "io/ioutil"
	//"context"
	"net/http"

	// config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	// "github.com/hypertrace/goagent/sdk/filter"
	// internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	// "github.com/hypertrace/goagent/sdk/internal/container"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// Server HTTP metrics.
const (
	// Pseudo of go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#RequestCount since a metric is not
	// created for that one for some reason.(annotated with hypertrace to avoid a duplicate if otel go ever implement
	// their own)
	RequestCount = "hypertrace.http.server.request_count" // Incoming request count total
)

type HttpOperationMetricsHandler struct {
	// operationName string
	operationNameGetter func(*http.Request) string
	// Some metrics in here.
	counters map[string]syncint64.Counter
}

var _ sdk.HttpOperationMetricsHandler = (*HttpOperationMetricsHandler)(nil)

// TODO: modify to return interface
func NewHttpOperationMetricsHandler(nameGetter func(*http.Request) string) sdk.HttpOperationMetricsHandler {
	return &HttpOperationMetricsHandler{
		operationNameGetter: nameGetter,
		counters:            make(map[string]syncint64.Counter, 1),
	}
}

func (mh *HttpOperationMetricsHandler) CreateRequestCount() {
	mp := global.MeterProvider()
	meter := mp.Meter("go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
		metric.WithInstrumentationVersion(otelhttp.SemVersion()))
	// counters := make(map[string]syncint64.Counter)

	requestCountCounter, err := meter.SyncInt64().Counter(RequestCount)
	if err != nil {
		otel.Handle(err)
	}

	mh.counters[RequestCount] = requestCountCounter
}

func (mh *HttpOperationMetricsHandler) AddToRequestCount(n int64, r *http.Request) {
	ctx := r.Context()
	labeler, _ := otelhttp.LabelerFromContext(ctx)
	operationName := mh.operationNameGetter(r)
	attributes := append(labeler.Get(), semconv.HTTPServerMetricAttributesFromHTTPRequest(operationName, r)...)
	mh.counters[RequestCount].Add(ctx, n, attributes...)
}

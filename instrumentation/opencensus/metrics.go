package opencensus // import "github.com/hypertrace/goagent/instrumentation/opencensus"

import (
	"net/http"

	"github.com/hypertrace/goagent/sdk"
)

// Will not support metrics for OC instrumentation
type httpOperationMetricsHandler struct {
}

func NewHttpOperationMetricsHandler() sdk.HttpOperationMetricsHandler {
	return &httpOperationMetricsHandler{}
}

func (mh *httpOperationMetricsHandler) AddToRequestCount(n int64, r *http.Request) {
}

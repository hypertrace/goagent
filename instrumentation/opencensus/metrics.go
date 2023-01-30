package opencensus // import "github.com/hypertrace/goagent/instrumentation/opencensus"

import (
	"net/http"

	"github.com/hypertrace/goagent/sdk"
)

type httpOperationMetricsHandler struct {
}

func NewHttpOperationMetricsHandler() sdk.HttpOperationMetricsHandler {
	return &httpOperationMetricsHandler{}
}

func (mh *httpOperationMetricsHandler) CreateRequestCount() {
}

func (mh *httpOperationMetricsHandler) AddToRequestCount(n int64, r *http.Request) {
}

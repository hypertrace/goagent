package sdk // import "github.com/hypertrace/goagent/sdk"

import (
	"net/http"
)

type HttpOperationMetricsHandler interface {
	CreateRequestCount()
	AddToRequestCount(int64, *http.Request)
}

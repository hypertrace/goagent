package hypergrpc // import "github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/hypergrpc"

import (
	"github.com/hypertrace/goagent/instrumentation/opencensus"
	"github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc/stats"
)

// WrapServerHandler wraps an OpenCensus ServerHandler and returns a new one that records
// the request/response body and metadata.
func WrapServerHandler(delegate *ocgrpc.ServerHandler) stats.Handler {
	return grpc.WrapStatsHandler(delegate, opencensus.SpanFromContext)
}

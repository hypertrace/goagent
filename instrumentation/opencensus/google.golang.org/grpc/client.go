package grpc

import (
	"github.com/traceableai/goagent/instrumentation/opencensus"
	traceablegrpc "github.com/traceableai/goagent/sdk/google.golang.org/grpc"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc/stats"
)

// WrapClientHandler wraps an OpenCensus ClientHandler and returns a new one that records
// the request/response body and metadata.
func WrapClientHandler(delegate *ocgrpc.ClientHandler) stats.Handler {
	return traceablegrpc.WrapStatsHandler(delegate, opencensus.SpanFromContext)
}

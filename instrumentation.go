package goagent

import (
	"net/http"

	"google.golang.org/grpc"
)

// Instrumentation defines the instrumentation elements for the APM.
// Every implementation is responsible to override them.
var Instrumentation = struct {
	// HttpHandler wraps a handler with instrumentation
	HTTPHandler     func(http.Handler) http.Handler
	GRPCInterceptor struct {
		UnaryClient  func() grpc.UnaryClientInterceptor
		StreamClient func() grpc.StreamClientInterceptor
		UnaryServer  func() grpc.UnaryServerInterceptor
		StreamServer func() grpc.StreamServerInterceptor
	}
}{
	// Default noop handler
	HTTPHandler: func(h http.Handler) http.Handler { return h },
}

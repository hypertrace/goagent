package hypergin // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gin-gonic/hypergin"

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

// Copied from: https://github.com/turtlemonvh/gin-wraphh
// In order to access the handler, we need to cast it to the nextRequestHandler type,
// which isn't exported if we use the package as a dep

type wrappedResponseWriter struct {
	gin.ResponseWriter
	writer http.ResponseWriter
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *wrappedResponseWriter) WriteString(s string) (n int, err error) {
	return w.writer.Write([]byte(s))
}

// An http.Handler that passes on calls to downstream middlewares
type nextRequestHandler struct {
	c *gin.Context
}

// Run the next request in the middleware chain and return
func (h *nextRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	savedCtx := h.c.Request.Context()
	defer func() {
		h.c.Request = h.c.Request.WithContext(savedCtx)
	}()

	h.c.Request = h.c.Request.WithContext(r.Context())
	h.c.Writer = &wrappedResponseWriter{h.c.Writer, w}
	h.c.Next()
}

// Wrap something that accepts an http.Handler, returns an http.Handler
func wrap(hh func(h http.Handler) http.Handler) gin.HandlerFunc {
	// Steps:
	// - create an http handler to pass `hh`
	// - call `hh` with the http handler, which returns a function
	// - call the ServeHTTP method of the resulting function to run the rest of the middleware chain
	return func(c *gin.Context) {
		// if we fail to extract the next request handler from delegate the route template won't be reported
		hh(&nextRequestHandler{c}).ServeHTTP(c.Writer, c.Request)

	}
}

// END OF copy

type ginRoute struct {
	route string
}

type hyperGinCtxKeyType string

const hyperGinKey hyperGinCtxKeyType = "gin_route"

func spanNameFormatter(operation string, r *http.Request) (spanName string) {
	routeWrapper, ok := r.Context().Value(hyperGinKey).(ginRoute)

	// if a ginRoute wasn't appended to the context we won't have the route path template(ex: /users/:id)
	// instead just report method
	if !ok {
		return r.Method
	}

	return routeWrapper.route
}

func Middleware(options *sdkhttp.Options) gin.HandlerFunc {
	return wrap(func(delegate http.Handler) http.Handler {
		wrappedHandler, ok := delegate.(*nextRequestHandler)
		ginOperationName := ""
		// if we fail to extract the next request handler from delegate the route template won't be reported
		if ok {
			ginOperationName := wrappedHandler.c.FullPath()
			rc := wrappedHandler.c.Request.Context()
			ctx := context.WithValue(rc, hyperGinKey, ginRoute{route: ginOperationName})
			wrappedHandler.c.Request = wrappedHandler.c.Request.WithContext(ctx)
		}
		return otelhttp.NewHandler(
			sdkhttp.WrapHandler(delegate, ginOperationName, opentelemetry.SpanFromContext, options, map[string]string{}),
			"",
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
		)
	})
}

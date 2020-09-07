package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/traceableai/goagent/instrumentation/internal"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

type handler struct {
	delegate          http.Handler
	defaultAttributes []label.KeyValue
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	if _, isNoop := span.(trace.NoopSpan); isNoop {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this Handler is not used
		// inside an instrumented Handler, hence we just invoke the delegate
		// round tripper.
		h.delegate.ServeHTTP(w, r)
		return
	}
	span.SetAttributes(h.defaultAttributes...)
	span.SetAttribute("http.url", r.URL.String())
	// Sets an attribute per each request header.
	setAttributesFromHeaders("request", r.Header, span)

	if shouldRecordBodyOfContentType(r.Header) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		// Only records the body if it is not empty and the content type
		// header is not streamable
		if len(body) > 0 {
			span.SetAttribute("http.request.body", string(body))

		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	// create http.ResponseWriter interceptor for tracking response size and
	// status code.
	wi := &rwInterceptor{w: w, statusCode: 200}

	// tag found response size and status code on exit
	defer func() {
		if len(wi.body) > 0 && shouldRecordBodyOfContentType(wi.Header()) {
			span.SetAttribute("http.response.body", string(wi.body))
		}

		// Sets an attribute per each response header.
		setAttributesFromHeaders("response", wi.Header(), span)

		span.End()
	}()

	h.delegate.ServeHTTP(wi, r)
}

// WrapHandler returns a new round tripper instrumented that relies on the
// needs to be used with OTel instrumentation.
func WrapHandler(delegate http.Handler) http.Handler {
	var defaultAttributes []label.KeyValue
	if containerID, err := internal.GetContainerID(); err != nil {
		defaultAttributes = append(defaultAttributes, label.String("container_id", containerID))
	}

	return &handler{delegate, defaultAttributes}
}

// Copied from Zipkin Go
// https://github.com/openzipkin/zipkin-go/blob/v0.2.3/middleware/http/server.go#L164
//
// rwInterceptor intercepts the ResponseWriter so it can track response size
// and returned status code.
type rwInterceptor struct {
	w          http.ResponseWriter
	body       []byte
	statusCode int
}

func (r *rwInterceptor) Header() http.Header {
	return r.w.Header()
}

func (r *rwInterceptor) Write(b []byte) (n int, err error) {
	n, err = r.w.Write(b)
	r.body = append(r.body, b...)
	return
}

func (r *rwInterceptor) WriteHeader(i int) {
	r.statusCode = i
	r.w.WriteHeader(i)
}

func (r *rwInterceptor) getStatusCode() int {
	return r.statusCode
}

func (r *rwInterceptor) wrap() http.ResponseWriter {
	var (
		hj, i0 = r.w.(http.Hijacker)
		cn, i1 = r.w.(http.CloseNotifier)
		pu, i2 = r.w.(http.Pusher)
		fl, i3 = r.w.(http.Flusher)
		rf, i4 = r.w.(io.ReaderFrom)
	)

	switch {
	case !i0 && !i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
		}{r}
	case !i0 && !i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
		}{r, rf}
	case !i0 && !i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Flusher
		}{r, fl}
	case !i0 && !i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
		}{r, fl, rf}
	case !i0 && !i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Pusher
		}{r, pu}
	case !i0 && !i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Pusher
			io.ReaderFrom
		}{r, pu, rf}
	case !i0 && !i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
		}{r, pu, fl}
	case !i0 && !i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{r, pu, fl, rf}
	case !i0 && i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
		}{r, cn}
	case !i0 && i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
		}{r, cn, rf}
	case !i0 && i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
		}{r, cn, fl}
	case !i0 && i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
			io.ReaderFrom
		}{r, cn, fl, rf}
	case !i0 && i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
		}{r, cn, pu}
	case !i0 && i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
			io.ReaderFrom
		}{r, cn, pu, rf}
	case !i0 && i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
			http.Flusher
		}{r, cn, pu, fl}
	case !i0 && i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{r, cn, pu, fl, rf}
	case i0 && !i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{r, hj}
	case i0 && !i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{r, hj, rf}
	case i0 && !i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
		}{r, hj, fl}
	case i0 && !i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
			io.ReaderFrom
		}{r, hj, fl, rf}
	case i0 && !i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
		}{r, hj, pu}
	case i0 && !i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			io.ReaderFrom
		}{r, hj, pu, rf}
	case i0 && !i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
		}{r, hj, pu, fl}
	case i0 && !i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{r, hj, pu, fl, rf}
	case i0 && i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
		}{r, hj, cn}
	case i0 && i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			io.ReaderFrom
		}{r, hj, cn, rf}
	case i0 && i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Flusher
		}{r, hj, cn, fl}
	case i0 && i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Flusher
			io.ReaderFrom
		}{r, hj, cn, fl, rf}
	case i0 && i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
		}{r, hj, cn, pu}
	case i0 && i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
			io.ReaderFrom
		}{r, hj, cn, pu, rf}
	case i0 && i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
			http.Flusher
		}{r, hj, cn, pu, fl}
	case i0 && i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{r, hj, cn, pu, fl, rf}
	default:
		return struct {
			http.ResponseWriter
		}{r}
	}
}

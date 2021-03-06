package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/container"
)

type handler struct {
	delegate                 http.Handler
	defaultAttributes        map[string]string
	spanFromContextRetriever sdk.SpanFromContext
	dataCaptureConfig        *config.DataCapture
	filter                   filter.Filter
}

// Options for HTTP handler instrumentation
type Options struct {
	Filter filter.Filter
}

// WrapHandler wraps an uninstrumented handler (e.g. a handleFunc) and returns a new one
// that should be used as base to an instrumented handler
func WrapHandler(delegate http.Handler, spanFromContext sdk.SpanFromContext, options *Options) http.Handler {
	defaultAttributes := make(map[string]string)
	if containerID, err := container.GetID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}
	var f filter.Filter = filter.NoOpFilter{}
	if options != nil && options.Filter != nil {
		f = options.Filter
	}
	return &handler{delegate, defaultAttributes, spanFromContext, internalconfig.GetConfig().GetDataCapture(), f}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := h.spanFromContextRetriever(r.Context())

	if span.IsNoop() {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this Handler is not used
		// inside an instrumented Handler, hence we just invoke the delegate
		// round tripper.
		h.delegate.ServeHTTP(w, r)
		return
	}

	for key, value := range h.defaultAttributes {
		span.SetAttribute(key, value)
	}

	url := r.URL.String()
	span.SetAttribute("http.url", url)

	// run filters on URL
	if h.filter.EvaluateURL(span, url) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	headers := r.Header
	// Sets an attribute per each request header.
	if h.dataCaptureConfig.HttpHeaders.Request.Value {
		SetAttributesFromHeaders("request", NewHeaderMapAccessor(r.Header), span)
	}

	// run filters on headers
	if h.filter.EvaluateHeaders(span, headers) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// nil check for body is important as this block turns the body into another
	// object that isn't nil and that will leverage the "Observer effect".
	if r.Body != nil && h.dataCaptureConfig.HttpBody.Request.Value && ShouldRecordBodyOfContentType(headerMapAccessor{r.Header}) {
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

		// run body filters
		if h.filter.EvaluateBody(span, body) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	// create http.ResponseWriter interceptor for tracking status code
	wi := &rwInterceptor{w: w, statusCode: 200}

	// tag found status code on exit
	defer func() {
		if h.dataCaptureConfig.HttpBody.Response.Value &&
			len(wi.body) > 0 &&
			ShouldRecordBodyOfContentType(headerMapAccessor{wi.Header()}) {
			span.SetAttribute("http.response.body", string(wi.body))
		}

		if h.dataCaptureConfig.HttpHeaders.Response.Value {
			// Sets an attribute per each response header.
			SetAttributesFromHeaders("response", headerMapAccessor{wi.Header()}, span)
		}
	}()

	h.delegate.ServeHTTP(wi, r)
}

// Copied from Zipkin Go
// https://github.com/openzipkin/zipkin-go/blob/v0.2.3/middleware/http/server.go#L164
//
// rwInterceptor intercepts the ResponseWriter so it can track returned status code.
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

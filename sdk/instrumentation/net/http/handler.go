package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	config "github.com/hypertrace/agent-config/gen/go/v1"
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
	mh                       sdk.HttpOperationMetricsHandler
}

// Options for HTTP handler instrumentation
type Options struct {
	Filter filter.Filter
}

// WrapHandler wraps an uninstrumented handler (e.g. a handleFunc) and returns a new one
// that should be used as base to an instrumented handler
func WrapHandler(delegate http.Handler, spanFromContext sdk.SpanFromContext, options *Options, spanAttributes map[string]string,
	mh sdk.HttpOperationMetricsHandler) http.Handler {
	defaultAttributes := make(map[string]string)
	for k, v := range spanAttributes {
		defaultAttributes[k] = v
	}
	if containerID, err := container.GetID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}
	var f filter.Filter = &filter.NoopFilter{}
	if options != nil && options.Filter != nil {
		f = options.Filter
	}

	return &handler{delegate, defaultAttributes, spanFromContext, internalconfig.GetConfig().GetDataCapture(), f, mh}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := h.spanFromContextRetriever(ctx)

	h.mh.AddToRequestCount(1, r)

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
	if strings.Contains(url, "://") {
		span.SetAttribute("http.url", url)
	} else {
		span.SetAttribute("http.target", url)
	}

	host := r.Host
	span.SetAttribute("http.request.header.host", host)

	// Sets an attribute per each request header.
	if h.dataCaptureConfig.HttpHeaders.Request.Value {
		SetAttributesFromHeaders("request", NewHeaderMapAccessor(r.Header), span)
	}

	// nil check for body is important as this block turns the body into another
	// object that isn't nil and that will leverage the "Observer effect".
	if r.Body != nil && h.dataCaptureConfig.HttpBody.Request.Value && ShouldRecordBodyOfContentType(headerMapAccessor{r.Header}) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		isMultipartFormDataBody := HasMultiPartFormDataContentTypeHeader(headerMapAccessor{r.Header})

		// Only records the body if it is not empty and the content type
		// header is not streamable
		if len(body) > 0 {
			setTruncatedBodyAttribute("request", body, int(h.dataCaptureConfig.BodyMaxSizeBytes.Value), span,
				isMultipartFormDataBody)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// single evaluation call to filter after capturing the configured parameters
	filterResult := h.filter.Evaluate(span)
	if filterResult.Block {
		w.WriteHeader(int(filterResult.ResponseStatusCode))
		return
	}

	// create http.ResponseWriter interceptor for tracking status code
	wi := &rwInterceptor{w: w, statusCode: 200}

	// tag found status code on exit
	defer func() {
		if h.dataCaptureConfig.HttpBody.Response.Value &&
			len(wi.body) > 0 &&
			ShouldRecordBodyOfContentType(headerMapAccessor{wi.Header()}) {
			setTruncatedBodyAttribute("response", wi.body, int(h.dataCaptureConfig.BodyMaxSizeBytes.Value), span,
				HasMultiPartFormDataContentTypeHeader(headerMapAccessor{wi.Header()}))
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

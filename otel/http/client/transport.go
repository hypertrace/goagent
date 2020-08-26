package client

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/traceableai/goagent/otel/http/internal"
	"go.opentelemetry.io/otel/api/trace"
)

var _ http.RoundTripper = &roundTripper{}

type roundTripper struct {
	delegate http.RoundTripper
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(req.Context())
	if s, isNoop := span.(trace.NoopSpan); isNoop {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this RoundTripper is not used
		// inside an instrumented transport, hence we just invoke the delegate
		// round tripper.
		return rt.delegate.RoundTrip(req)
	}

	for key, value := range req.Header {
		span.SetAttribute("http.request.header."+key, value)
	}

	// Only records the body if it is not empty and the content type header
	// is in the recording accept list. Notice in here we rely on the fact that
	// the content type is not streamable, otherwise we could end up in a very
	// expensive parsing of a big body in memory.
	if internal.IsContentTypeInAllowList(req.Header) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return rt.delegate.RoundTrip(req)
		}
		defer req.Body.Close()

		if len(body) > 0 {
			span.SetAttribute("http.request.body", string(body))
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	res, err := rt.delegate.RoundTrip(req)

	// Notice, parsing a streamed content in memory can be expensive.
	if internal.IsContentTypeInAllowList(res.Header) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return rt.delegate.RoundTrip(req)
		}
		defer req.Body.Close()

		if len(body) > 0 {
			span.SetAttribute("http.response.body", string(body))
		}

		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	// Sets an attribute per each response header.
	for key, value := range res.Header {
		span.SetAttribute("http.response.header."+key, value)
	}

	return res, err
}

// Wrap returns a new transport instrumented by OTel and
// records body and headers.
func Wrap(delegate http.RoundTripper) http.RoundTripper {
	return &roundTripper{delegate}
}

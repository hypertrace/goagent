package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
)

var _ http.RoundTripper = &roundTripper{}

type roundTripper struct {
	delegate http.RoundTripper
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(req.Context())
	if _, isNoop := span.(trace.NoopSpan); isNoop {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this RoundTripper is not used
		// inside an instrumented transport, hence we just invoke the delegate
		// round tripper.
		return rt.delegate.RoundTrip(req)
	}

	for key, value := range req.Header {
		span.SetAttribute("http.request.header."+key, value[0])
	}

	// Only records the body if it is not empty and the content type header
	// is in the recording accept list. Notice in here we rely on the fact that
	// the content type is not streamable, otherwise we could end up in a very
	// expensive parsing of a big body in memory.
	if shouldRecordBodyOfContentType(req.Header) {
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
	if err != nil {
		return res, err
	}

	// Notice, parsing a streamed content in memory can be expensive.
	if shouldRecordBodyOfContentType(res.Header) {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return res, nil
		}
		defer res.Body.Close()

		if len(body) > 0 {
			span.SetAttribute("http.response.body", string(body))
		}

		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	// Sets an attribute per each response header.
	for key, value := range res.Header {
		span.SetAttribute("http.response.header."+key, value[0])
	}

	return res, err
}

// WrapTransport returns a new round tripper instrumented that relies on the
// needs to be used with OTel instrumentation.
func WrapTransport(delegate http.RoundTripper) http.RoundTripper {
	return &roundTripper{delegate}
}

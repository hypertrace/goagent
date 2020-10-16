package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/traceableai/goagent/sdk"
	"github.com/traceableai/goagent/sdk/internal"
)

var _ http.RoundTripper = &roundTripper{}

type roundTripper struct {
	delegate                 http.RoundTripper
	defaultAttributes        map[string]string
	spanFromContextRetriever sdk.SpanFromContext
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	span := rt.spanFromContextRetriever(req.Context())
	if span.IsNoop() {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this RoundTripper is not used
		// inside an instrumented transport, hence we just invoke the delegate
		// round tripper.
		return rt.delegate.RoundTrip(req)
	}

	for key, value := range rt.defaultAttributes {
		span.SetAttribute(key, value)
	}

	setAttributesFromHeaders("request", req.Header, span)

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
	setAttributesFromHeaders("response", res.Header, span)

	return res, err
}

// WrapTransport returns a new http.RoundTripper that should be wrapped
// by an instrumented http.RoundTripper
func WrapTransport(delegate http.RoundTripper, spanFromContextRetriever sdk.SpanFromContext) http.RoundTripper {
	defaultAttributes := make(map[string]string)
	if containerID, err := internal.GetContainerID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	return &roundTripper{delegate, defaultAttributes, spanFromContextRetriever}
}

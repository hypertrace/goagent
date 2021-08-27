package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/container"
)

var _ http.RoundTripper = &roundTripper{}

type roundTripper struct {
	delegate                 http.RoundTripper
	defaultAttributes        map[string]string
	spanFromContextRetriever sdk.SpanFromContext
	dataCaptureConfig        *config.DataCapture
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

	if rt.dataCaptureConfig.HttpHeaders.Request.Value {
		SetAttributesFromHeaders("request", headerMapAccessor{req.Header}, span)
	}

	// Only records the body if it is not empty and the content type header
	// is in the recording accept list. Notice in here we rely on the fact that
	// the content type is not streamable, otherwise we could end up in a very
	// expensive parsing of a big body in memory.
	if req.Body != nil && rt.dataCaptureConfig.HttpBody.Request.Value && ShouldRecordBodyOfContentType(headerMapAccessor{req.Header}) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return rt.delegate.RoundTrip(req)
		}
		defer req.Body.Close()

		if len(body) > 0 {
			setTruncatedBodyAttribute("request", body, int(rt.dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	res, err := rt.delegate.RoundTrip(req)
	if err != nil {
		return res, err
	}

	// Notice, parsing a streamed content in memory can be expensive.
	if rt.dataCaptureConfig.HttpBody.Response.Value && ShouldRecordBodyOfContentType(headerMapAccessor{res.Header}) {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return res, nil
		}
		defer res.Body.Close()

		if len(body) > 0 {
			setTruncatedBodyAttribute("response", body, int(rt.dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		}

		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	if rt.dataCaptureConfig.HttpHeaders.Response.Value {
		// Sets an attribute per each response header.
		SetAttributesFromHeaders("response", headerMapAccessor{res.Header}, span)
	}

	return res, err
}

// WrapTransport returns a new http.RoundTripper that should be wrapped
// by an instrumented http.RoundTripper
func WrapTransport(delegate http.RoundTripper, spanFromContextRetriever sdk.SpanFromContext) http.RoundTripper {
	defaultAttributes := make(map[string]string)
	if containerID, err := container.GetID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	return &roundTripper{delegate, defaultAttributes, spanFromContextRetriever, internalconfig.GetConfig().GetDataCapture()}
}

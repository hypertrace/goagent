package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	codes "github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/container"
)

var _ http.RoundTripper = &roundTripper{}

type roundTripper struct {
	delegate                 http.RoundTripper
	defaultAttributes        map[string]string
	spanFromContextRetriever sdk.SpanFromContext
	dataCaptureConfig        *config.DataCapture
	filter                   filter.Filter
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
	reqHeadersAccessor := NewHeaderMapAccessor(req.Header)

	for key, value := range rt.defaultAttributes {
		span.SetAttribute(key, value)
	}

	if rt.dataCaptureConfig.HttpHeaders.Request.Value {
		SetAttributesFromHeaders("request", reqHeadersAccessor, span)
	}

	// Only records the body if it is not empty and the content type header
	// is in the recording accept list. Notice in here we rely on the fact that
	// the content type is not streamable, otherwise we could end up in a very
	// expensive parsing of a big body in memory.
	if req.Body != nil && rt.dataCaptureConfig.HttpBody.Request.Value && ShouldRecordBodyOfContentType(reqHeadersAccessor) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return rt.delegate.RoundTrip(req)
		}
		defer req.Body.Close()

		if len(body) > 0 {
			setTruncatedBodyAttribute("request", body, int(rt.dataCaptureConfig.BodyMaxSizeBytes.Value), span,
				HasMultiPartFormDataContentTypeHeader(reqHeadersAccessor))
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	filterResult := rt.filter.Evaluate(span)
	if filterResult.Block {
		span.SetStatus(codes.StatusCodeError, "Access Denied")
		span.SetAttribute("http.status_code", filterResult.ResponseStatusCode)
		return &http.Response{
			Status:     "Access Denied",
			StatusCode: int(filterResult.ResponseStatusCode),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Request:    req,
			Header: map[string][]string{
				"Content-Type": {"text/plain"},
			},
			Body: io.NopCloser(strings.NewReader(filterResult.ResponseMessage)),
		}, nil
	} else if filterResult.Decorations != nil {
		for _, header := range filterResult.Decorations.RequestHeaderInjections {
			req.Header.Add(header.Key, header.Value)
			span.SetAttribute("http.request.header."+header.Key, header.Value)
		}
	}

	res, err := rt.delegate.RoundTrip(req)
	if err != nil {
		return res, err
	}
	resHeadersAccessor := NewHeaderMapAccessor(res.Header)

	// Notice, parsing a streamed content in memory can be expensive.
	if rt.dataCaptureConfig.HttpBody.Response.Value && ShouldRecordBodyOfContentType(resHeadersAccessor) {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return res, nil
		}
		defer res.Body.Close()

		if len(body) > 0 {
			setTruncatedBodyAttribute("response", body, int(rt.dataCaptureConfig.BodyMaxSizeBytes.Value), span,
				HasMultiPartFormDataContentTypeHeader(resHeadersAccessor))
		}

		res.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	if rt.dataCaptureConfig.HttpHeaders.Response.Value {
		// Sets an attribute per each response header.
		SetAttributesFromHeaders("response", resHeadersAccessor, span)
	}

	return res, err
}

// WrapTransport returns a new http.RoundTripper that should be wrapped
// by an instrumented http.RoundTripper
func WrapTransport(delegate http.RoundTripper, spanFromContextRetriever sdk.SpanFromContext, options *Options, spanAttributes map[string]string) http.RoundTripper {
	defaultAttributes := make(map[string]string)
	for k, v := range spanAttributes {
		defaultAttributes[k] = v
	}
	if containerID, err := container.GetID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	var filter filter.Filter = &filter.NoopFilter{}
	if options != nil && options.Filter != nil {
		filter = options.Filter
	}
	return &roundTripper{
		delegate:                 delegate,
		defaultAttributes:        defaultAttributes,
		spanFromContextRetriever: spanFromContextRetriever,
		dataCaptureConfig:        internalconfig.GetConfig().GetDataCapture(),
		filter:                   filter,
	}
}

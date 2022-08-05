package grpc

import (
	"context"
	"testing"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter"
	"github.com/hypertrace/goagent/sdk/filterutils"
	"github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc/internal/helloworld"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func makeMockUnaryServerInterceptor(mockSpans *[]*mock.Span) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span := mock.NewSpan()
		ctx = mock.ContextWithSpan(ctx, span)
		*mockSpans = append(*mockSpans, span)
		return handler(ctx, req)
	}
}

func TestServerInterceptorHelloWorldSuccess(t *testing.T) {
	defer internalconfig.ResetConfig()

	spans := []*mock.Span{}
	mockUnaryInterceptor := makeMockUnaryServerInterceptor(&spans)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			WrapUnaryServerInterceptor(mockUnaryInterceptor, mock.SpanFromContext, &Options{}, map[string]string{"foo": "bar"}),
		),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{})

	dialer := createDialer(s)

	ctx := context.Background()

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key", "test_value"))
	_, err = client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Pupo",
	})
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	assert.Equal(t, 1, len(spans))

	span := spans[0]

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value", span.ReadAttribute("rpc.request.metadata.test_key").(string))
	assert.Equal(t, "POST", span.ReadAttribute("rpc.request.metadata.:method").(string))
	assert.Equal(t, "bar", span.ReadAttribute("foo").(string))

	expectedBody := "{\"name\":\"Pupo\"}"
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Pupo\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerInterceptorFilter(t *testing.T) {
	defer internalconfig.ResetConfig()

	tCases := map[string]struct {
		expectedFilterResult bool
		multiFilter          *filter.MultiFilter
	}{
		"no filter": {
			expectedFilterResult: false,
			multiFilter:          filter.NewMultiFilter(),
		},
		"headers filter": {
			expectedFilterResult: true,
			multiFilter: filter.NewMultiFilter(mock.Filter{
				URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) filterutils.FilterResult {
					assert.Equal(t, []string{"test_value"}, headers["test_key"])
					return filterutils.FilterResult{Block: true, ResponseStatusCode: 403}
				},
			}),
		},
		"body filter": {
			expectedFilterResult: true,
			multiFilter: filter.NewMultiFilter(mock.Filter{
				BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) filterutils.FilterResult {
					assert.Equal(t, "{\"name\":\"Pupo\"}", string(body))
					return filterutils.FilterResult{Block: true, ResponseStatusCode: 403}
				},
			}),
		},
	}

	spans := []*mock.Span{}
	mockUnaryInterceptor := makeMockUnaryServerInterceptor(&spans)

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			// wrap interceptor with filter
			s := grpc.NewServer(
				grpc.UnaryInterceptor(
					WrapUnaryServerInterceptor(mockUnaryInterceptor, mock.SpanFromContext, &Options{Filter: tCase.multiFilter}, map[string]string{}),
				),
			)
			defer s.Stop()

			helloworld.RegisterGreeterServer(s, &server{})

			dialer := createDialer(s)

			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"bufnet",
				grpc.WithContextDialer(dialer),
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				t.Fatalf("failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			client := helloworld.NewGreeterClient(conn)

			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key", "test_value"))
			_, err = client.SayHello(ctx, &helloworld.HelloRequest{
				Name: "Pupo",
			})
			if tCase.expectedFilterResult {
				assert.Equal(t, codes.Code(403), status.Code(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestServerInterceptorFilterWithMaxProcessingBodyLen(t *testing.T) {
	spans := []*mock.Span{}
	mockUnaryInterceptor := makeMockUnaryServerInterceptor(&spans)

	cfg := &config.AgentConfig{
		DataCapture: &config.DataCapture{
			RpcBody: &config.Message{
				Request: config.Bool(true),
			},
			BodyMaxProcessingSizeBytes: config.Int32(1),
		},
	}
	cfg.LoadFromEnv()

	internalconfig.InitConfig(cfg)
	defer internalconfig.ResetConfig()

	// wrap interceptor with filter
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			WrapUnaryServerInterceptor(mockUnaryInterceptor, mock.SpanFromContext, &Options{Filter: mock.Filter{
				BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) filterutils.FilterResult {
					assert.Equal(t, "{", string(body)) // body is truncated
					return filterutils.FilterResult{}
				},
			}}, nil),
		),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key", "test_value"))
	_, err = client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Pupo",
	})

	assert.NoError(t, err)
}

func TestServerHandlerHelloWorldSuccess(t *testing.T) {
	defer internalconfig.ResetConfig()

	mockHandler := &mockHandler{}

	s := grpc.NewServer(
		grpc.StatsHandler(WrapStatsHandler(mockHandler, mock.SpanFromContext)),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUserAgent("test_agent"),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key", "test_value"))

	_, err = client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Pupo",
	})
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	assert.Equal(t, 1, len(mockHandler.Spans))

	span := mockHandler.Spans[0]

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value", span.ReadAttribute("rpc.request.metadata.test_key").(string))

	assert.Equal(t, "bufnet", span.ReadAttribute("rpc.request.metadata.:authority").(string))
	assert.Equal(t, "application/grpc", span.ReadAttribute("rpc.request.metadata.content-type").(string))
	assert.Contains(t, span.ReadAttribute("rpc.request.metadata.user-agent").(string), "test_agent")

	expectedBody := "{\"name\":\"Pupo\"}"
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Pupo\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}

type fakeALTSAuthInfo struct{}

func (fakeALTSAuthInfo) AuthType() string {
	return "tls"
}

func TestSetSchemeAttributes(t *testing.T) {
	tCases := map[string]struct {
		expectedScheme string
		AuthInfo       credentials.AuthInfo
	}{
		"no auth info": {
			expectedScheme: "http",
			AuthInfo:       nil,
		},
		"with auth info": {
			expectedScheme: "https",
			AuthInfo:       fakeALTSAuthInfo{},
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			ms := mock.NewSpan()

			p := &peer.Peer{
				Addr:     nil,
				AuthInfo: tCase.AuthInfo,
			}

			pctx := peer.NewContext(ctx, p)
			setSchemeAttributes(pctx, ms)
			assert.Equal(t, tCase.expectedScheme, ms.ReadAttribute("rpc.request.metadata.:scheme"))
		})
	}
}

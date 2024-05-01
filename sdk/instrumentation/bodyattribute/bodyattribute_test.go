package bodyattribute

import (
	"encoding/base64"
	"testing"

	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestBodyTruncationSuccess(t *testing.T) {
	s := mock.NewSpan()
	SetTruncatedBodyAttribute("http.request.body", []byte("text"), 2, s)
	assert.Equal(t, "te", s.ReadAttribute("http.request.body"))
	assert.True(t, (s.ReadAttribute("http.request.body.truncated")).(bool))
	assert.Zero(t, s.RemainingAttributes())
}

func TestBodyTruncationIsSkipped(t *testing.T) {
	s := mock.NewSpan()
	SetTruncatedBodyAttribute("rpc.response.body", []byte("text"), 7, s)
	assert.Equal(t, "text", s.ReadAttribute("rpc.response.body"))
	assert.Zero(t, s.RemainingAttributes())
}

func TestBodyTruncationEmptyBody(t *testing.T) {
	s := mock.NewSpan()
	SetTruncatedBodyAttribute("body_attr", []byte{}, 7, s)
	assert.Nil(t, s.ReadAttribute("body_attr"))
	assert.Zero(t, s.RemainingAttributes())
}

func TestSetTruncatedEncodedBodyAttributeNoTruncation(t *testing.T) {
	s := mock.NewSpan()
	SetTruncatedEncodedBodyAttribute("http.request.body", []byte("text"), 7, s)
	assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte("text")), s.ReadAttribute("http.request.body.base64"))
	assert.Zero(t, s.RemainingAttributes())
}

func TestSetTruncatedEncodedBodyAttribute(t *testing.T) {
	s := mock.NewSpan()
	SetTruncatedEncodedBodyAttribute("http.request.body", []byte("text"), 2, s)
	assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte("te")), s.ReadAttribute("http.request.body.base64"))
	assert.True(t, (s.ReadAttribute("http.request.body.truncated")).(bool))
	assert.Zero(t, s.RemainingAttributes())
}

func TestSetTruncatedEncodedBodyAttributeEmptyBody(t *testing.T) {
	s := mock.NewSpan()
	SetTruncatedEncodedBodyAttribute("http.request.body", []byte{}, 2, s)
	assert.Nil(t, s.ReadAttribute("http.request.body.base64"))
	assert.Zero(t, s.RemainingAttributes())
}

func TestSetBodyAttribute(t *testing.T) {
	testBody := "test1test2"
	type args struct {
		attrName  string
		body      []byte
		truncated bool
		span      *mock.Span
	}
	tests := []struct {
		name               string
		args               args
		expectedAssertions func(t *testing.T, gotSpan *mock.Span)
	}{
		{
			name: "empty body, truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(""),
				truncated: true,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Nil(t, gotSpan.ReadAttribute("http.request.body"))
				assert.Zero(t, gotSpan.RemainingAttributes())
			},
		},
		{
			name: "empty body, not truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(""),
				truncated: false,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Nil(t, gotSpan.ReadAttribute("http.request.body"))
				assert.Zero(t, gotSpan.RemainingAttributes())
			},
		},
		{
			name: "non empty body, not truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(testBody),
				truncated: false,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Equal(t, testBody, gotSpan.ReadAttribute("http.request.body"))
				assert.Zero(t, gotSpan.RemainingAttributes())
			},
		},
		{
			name: "non empty body, truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(testBody),
				truncated: true,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Equal(t, testBody, gotSpan.ReadAttribute("http.request.body"))
				assert.True(t, (gotSpan.ReadAttribute("http.request.body.truncated")).(bool))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetBodyAttribute(tt.args.attrName, tt.args.body, tt.args.truncated, tt.args.span)
			tt.expectedAssertions(t, tt.args.span)
		})
	}
}

func TestSetEncodedBodyAttribute(t *testing.T) {
	testBody := "test1test2"
	type args struct {
		attrName  string
		body      []byte
		truncated bool
		span      *mock.Span
	}
	tests := []struct {
		name               string
		args               args
		expectedAssertions func(t *testing.T, gotSpan *mock.Span)
	}{
		{
			name: "empty body, truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(""),
				truncated: true,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Nil(t, gotSpan.ReadAttribute("http.request.body.base64"))
				assert.Zero(t, gotSpan.RemainingAttributes())
			},
		},
		{
			name: "empty body, not truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(""),
				truncated: false,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Nil(t, gotSpan.ReadAttribute("http.request.body.base64"))
				assert.Zero(t, gotSpan.RemainingAttributes())
			},
		},
		{
			name: "non empty body, not truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(testBody),
				truncated: false,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte(testBody)), gotSpan.ReadAttribute("http.request.body.base64"))
				assert.Zero(t, gotSpan.RemainingAttributes())
			},
		},
		{
			name: "non empty body, truncated",
			args: args{
				attrName:  "http.request.body",
				body:      []byte(testBody),
				truncated: true,
				span:      mock.NewSpan(),
			},
			expectedAssertions: func(t *testing.T, gotSpan *mock.Span) {
				assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte(testBody)), gotSpan.ReadAttribute("http.request.body.base64"))
				assert.True(t, (gotSpan.ReadAttribute("http.request.body.truncated")).(bool))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetEncodedBodyAttribute(tt.args.attrName, tt.args.body, tt.args.truncated, tt.args.span)
			tt.expectedAssertions(t, tt.args.span)
		})
	}
}

func TestSetBodyWithoutUtf8(t *testing.T) {
	multiByteCharString := []byte("こんにちは世界こんにちは世界こんにちは世界こんにちは世界こんにちは世界")
	span := mock.NewSpan()
	SetTruncatedBodyAttribute("http.request.body", multiByteCharString, 23, span)
	value := span.ReadAttribute("http.request.body")
	assert.Equal(t, value.(string), "こんにちは世界")
	v := len(value.(string))
	assert.Equal(t, v, 21)
}

func TestSetB64BodyWithoutUtf8(t *testing.T) {
	multiByteCharString := []byte("こんにちは世界こんにちは世界こんにちは世界こんにちは世界こんにちは世界")
	span := mock.NewSpan()
	SetTruncatedEncodedBodyAttribute("http.request.body", multiByteCharString, 23, span)
	value := span.ReadAttribute("http.request.body.base64")
	decodedBytes, err := base64.StdEncoding.DecodeString(value.(string))
	assert.NoError(t, err)
	assert.Equal(t, string(decodedBytes), "こんにちは世界")
	v := len(decodedBytes)
	assert.Equal(t, v, 21)
}

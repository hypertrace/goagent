package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIsSuccessfullyTraced(t *testing.T) {
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("test_response_body"))
	})

	ih := NewHandler(h)

	r, _ := http.NewRequest("GET", "http://traceable.ai", strings.NewReader("test_request_body"))
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	if want, have := "test_response_body", w.Body.String(); want != have {
		t.Errorf("unexpected response body, want %q, have %q", want, have)
	}
}

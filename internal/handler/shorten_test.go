package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devkyudin/shortener/internal/service"
)

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func TestShorten(t *testing.T) {
	type testCase struct {
		name            string
		req             *http.Request
		wantStatus      int
		wantBody        string
		wantContentType string
	}

	tests := []testCase{
		{
			name: "returns 201 and short link for valid POST text/plain",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
				r.Header.Set("Content-Type", "text/plain")
				return r
			}(),
			wantStatus:      http.StatusCreated,
			wantBody:        service.CreateShortLink("https://example.com"),
			wantContentType: "text/plain",
		},
		{
			name: "returns 400 for non-POST method",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("https://example.com"))
				r.Header.Set("Content-Type", "text/plain")
				return r
			}(),
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
		{
			name: "returns 400 for wrong content type",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
				r.Header.Set("Content-Type", "application/json")
				return r
			}(),
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
		{
			name: "returns 400 for empty body",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				r.Header.Set("Content-Type", "text/plain")
				return r
			}(),
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
		{
			name: "returns 400 when body read errors",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", errReader{})
				r.Header.Set("Content-Type", "text/plain")
				return r
			}(),
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			Shorten(rec, tt.req)

			res := rec.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantContentType != "" {
				if res.Header.Get("Content-Type") != tt.wantContentType {
					t.Fatalf("Content-Type: got %q, want %q", res.Header.Get("Content-Type"), tt.wantContentType)
				}
			}
			if string(body) != tt.wantBody {
				t.Fatalf("body: got %q, want %q", string(body), tt.wantBody)
			}
		})
	}
}

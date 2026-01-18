package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/service"
)

func TestGetLink(t *testing.T) {
	type testCase struct {
		name            string
		req             *http.Request
		wantStatus      int
		wantLocation    string
		wantContentType string
	}

	tests := []testCase{
		{
			name: "returns 307 and Location for existing id",
			req: func() *http.Request {
				shortLink := service.CreateShortLink("https://example.com")
				id := shortLink[len("http://localhost:8080/"):]
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.SetPathValue("id", id)
				return r
			}(),
			wantStatus:      http.StatusTemporaryRedirect,
			wantLocation:    "https://example.com",
			wantContentType: "text/plain",
		},
		{
			name: "returns 400 for non-GET method",
			req: func() *http.Request {
				shortLink := service.CreateShortLink("https://example.com")
				id := shortLink[len("http://localhost:8080/"):]
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.SetPathValue("id", id)
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "returns 400 for missing id",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "returns 400 for unknown id",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/unknown-id", nil)
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			GetLink(rec, tt.req)

			res := rec.Result()
			defer res.Body.Close()
			_, _ = io.ReadAll(res.Body)

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if tt.wantContentType != "" {
				if res.Header.Get("Content-Type") != tt.wantContentType {
					t.Fatalf("Content-Type: got %q, want %q", res.Header.Get("Content-Type"), tt.wantContentType)
				}
			}
			if tt.wantLocation != "" {
				if res.Header.Get("Location") != tt.wantLocation {
					t.Fatalf("Location: got %q, want %q", res.Header.Get("Location"), tt.wantLocation)
				}
			}
		})
	}
}

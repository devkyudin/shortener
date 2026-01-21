package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetLink(t *testing.T) {
	type testCase struct {
		name string
		req  req
		want want
	}

	tests := []testCase{
		{
			name: "returns 307 and Location for existing id",
			req: req{
				url:         "/" + service.CreateShortLink("https://example.com")[len("http://localhost:8080/"):],
				methodName:  "GET",
				contentType: "text/plain",
			},
			want: want{
				status:      http.StatusTemporaryRedirect,
				location:    "https://example.com",
				contentType: "text/plain",
			},
		},
		{
			name: "returns 400 for non-GET method",
			req: req{
				url:         "/" + service.CreateShortLink("https://example.com")[len("http://localhost:8080/"):],
				methodName:  "POST",
				contentType: "text/plain",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "returns 400 for missing id",
			req: req{
				url:         "/",
				methodName:  "GET",
				contentType: "text/plain",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "returns 400 for unknown id",
			req: req{
				url:         "/unknown-id",
				methodName:  "GET",
				contentType: "text/plain",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}

	ts := httptest.NewServer(ShortenerRouter())
	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testRequest(t, httptest.NewServer(ShortenerRouter()), tt.req)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			}(res.Body)
			assert.Equal(t, tt.want.status, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

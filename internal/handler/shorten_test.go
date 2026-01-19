package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func TestShorten(t *testing.T) {
	type testCase struct {
		name string
		req  req
		want want
	}

	tests := []testCase{
		{
			name: "returns 201 and short link for valid POST text/plain",
			req: req{
				url:         "/",
				methodName:  "POST",
				contentType: "text/plain",
				body:        "https://example.com",
			},
			want: want{
				status:      http.StatusCreated,
				body:        service.CreateShortLink("https://example.com"),
				contentType: "text/plain",
			},
		},
		{
			name: "returns 400 for non-POST method",
			req: req{
				url:         "/",
				methodName:  "GET",
				contentType: "text/plain",
				body:        "https://example.com",
			},
			want: want{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "returns 400 for wrong content type",
			req: req{
				url:         "/",
				methodName:  "POST",
				contentType: "application/json",
				body:        "https://example.com",
			},
			want: want{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "returns 400 for empty body",
			req: req{
				url:         "/",
				methodName:  "POST",
				contentType: "text/plain",
				body:        "",
			},
			want: want{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testRequest(t, httptest.NewServer(ShortenerRouter()), tt.req)
			assert.Equal(t, tt.want.status, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

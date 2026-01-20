package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/devkyudin/shortener/internal/service"
	"github.com/devkyudin/shortener/internal/testUtils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type req struct {
	url         string
	methodName  string
	contentType string
	body        string
}

type want struct {
	status      int
	body        string
	contentType string
	location    string
}

func TestMain(m *testing.M) {
	testUtils.SetupTestEnvironment()
	code := m.Run()
	os.Exit(code)
}

func TestShortenerRouter(t *testing.T) {
	tests := []struct {
		name string
		req  req
		want want
	}{
		{
			name: "Shorten: Bad request empty body",
			req: req{
				url:         "/",
				methodName:  "POST",
				contentType: "text/plain",
				body:        "",
			},
			want: want{
				status: 400,
				body:   "",
			},
		},
		{
			name: "Shorten: Bad request wrong content type",
			req: req{
				url:         "/",
				methodName:  "POST",
				contentType: "application/json",
				body:        "https://example.com",
			},
			want: want{
				status: 400,
				body:   "",
			},
		},
		{
			name: "Shorten: Bad request wrong method",
			req: req{
				url:         "/",
				methodName:  "GET",
				contentType: "text/plain",
				body:        "https://example.com",
			},
			want: want{
				status: 400,
				body:   "",
			},
		},
		{
			name: "Shorten: Created short link",
			req: req{
				url:         "/",
				methodName:  "POST",
				contentType: "text/plain",
				body:        "https://example.com",
			},
			want: want{
				status: 201,
				body:   service.CreateShortLink("https://example.com"),
			},
		},
		{
			name: "GetLink: Bad request missing id",
			req: req{
				url:         "/",
				methodName:  "GET",
				contentType: "text/plain",
				body:        "",
			},
			want: want{
				status: 400,
				body:   "",
			},
		},
		{
			name: "GetLink: Bad request unknown id",
			req: req{
				url:         "/unknown-id",
				methodName:  "GET",
				contentType: "text/plain",
				body:        "",
			},
			want: want{
				status: 400,
				body:   "",
			},
		},
		{
			name: "GetLink: Bad request wrong method",
			req: req{
				url:         "/some-id",
				methodName:  "POST",
				contentType: "text/plain",
				body:        "",
			},
			want: want{
				status: 400,
				body:   "",
			},
		},
		{
			name: "GetLink: Temporary redirect for existing id",
			req: req{
				url:         "/" + service.CreateShortLink("https://example.com")[len("http://localhost:8080/"):],
				methodName:  "GET",
				contentType: "text/plain",
				body:        "",
			},
			want: want{
				status: 307,
				body:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testRequest(t, httptest.NewServer(ShortenerRouter()), tt.req)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.status, resp.StatusCode)
			assert.Equal(t, tt.want.body, get)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, req req) (*http.Response, string) {
	reader := io.NopCloser(io.Reader(strings.NewReader(req.body)))
	r, err := http.NewRequest(req.methodName, ts.URL+req.url, reader)
	require.NoError(t, err)

	r.Header.Set("Content-Type", req.contentType)
	res, err := ts.Client().Do(r)
	require.NoError(t, err)
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res, string(resp)
}

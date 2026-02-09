package router_test

import (
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/dependencies"
	"github.com/devkyudin/shortener/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestShortenerRouter(t *testing.T) {
	deps := dependencies.GetDependencies()
	tests := []struct {
		name string
		req  testutils.Req
		want testutils.Want
	}{
		{
			name: "Shorten JSON: Bad request empty body",
			req: testutils.Req{
				URL:         "/api/shorten",
				MethodName:  "POST",
				ContentType: "application/json",
				Body:        "{}",
			},
			want: testutils.Want{
				Status: 400,
				Body:   "",
			},
		},
		{
			name: "Shorten JSON: Bad request wrong content type",
			req: testutils.Req{
				URL:         "/api/shorten",
				MethodName:  "POST",
				ContentType: "text/plain",
				Body:        `{"url":"https://example.com"}`,
			},
			want: testutils.Want{
				Status: 400,
				Body:   "",
			},
		},
		{
			name: "Shorten JSON: Created short link",
			req: testutils.Req{
				URL:         "/api/shorten",
				MethodName:  "POST",
				ContentType: "application/json",
				Body:        `{"url":"https://some-new-example.com"}`,
			},
			want: testutils.Want{
				Status: 201,
				Body:   `{"result":"` + testutils.CreateShortLinkSafe(deps.URLService, "https://some-new-example.com") + `"}` + "\n",
			},
		},
		{
			name: "Shorten PlainText: Bad request empty body",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "POST",
				ContentType: "text/plain",
				Body:        "",
			},
			want: testutils.Want{
				Status: 400,
				Body:   "",
			},
		},
		{
			name: "Shorten PlainText: Bad request wrong content type",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "POST",
				ContentType: "application/json",
				Body:        "https://example.com",
			},
			want: testutils.Want{
				Status: 400,
				Body:   "",
			},
		},
		{
			name: "Shorten PlainText: Bad request wrong method",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "GET",
				ContentType: "text/plain",
				Body:        "https://example.com",
			},
			want: testutils.Want{
				Status: 405,
				Body:   "",
			},
		},
		{
			name: "Shorten PlainText: Created short link",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "POST",
				ContentType: "text/plain",
				Body:        "https://example.com",
			},
			want: testutils.Want{
				Status: 201,
				Body:   testutils.CreateShortLinkSafe(deps.URLService, "https://example.com"),
			},
		},
		{
			name: "GetLink: Bad request missing id",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "GET",
				ContentType: "text/plain",
				Body:        "",
			},
			want: testutils.Want{
				Status: 405,
				Body:   "",
			},
		},
		{
			name: "GetLink: Bad request unknown id",
			req: testutils.Req{
				URL:         "/unknown-id",
				MethodName:  "GET",
				ContentType: "text/plain",
				Body:        "",
			},
			want: testutils.Want{
				Status: 400,
				Body:   "",
			},
		},
		{
			name: "GetLink: Bad request wrong method",
			req: testutils.Req{
				URL:         "/some-id",
				MethodName:  "POST",
				ContentType: "text/plain",
				Body:        "",
			},
			want: testutils.Want{
				Status: 405,
				Body:   "",
			},
		},
		{
			name: "GetLink: Temporary redirect for existing id",
			req: testutils.Req{
				URL:         "/" + testutils.CreateShortLinkSafe(deps.URLService, "https://example.com")[len("http://localhost:8080/"):],
				MethodName:  "GET",
				ContentType: "text/plain",
				Body:        "",
			},
			want: testutils.Want{
				Status: 307,
				Body:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, get := testutils.TestRequest(t, httptest.NewServer(*deps.Router), tt.req)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.Status, resp.StatusCode)
			assert.Equal(t, tt.want.Body, get)
		})
	}
}

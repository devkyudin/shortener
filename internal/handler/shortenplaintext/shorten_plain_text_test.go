package shortenplaintext_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/dependencies"
	"github.com/devkyudin/shortener/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	type testCase struct {
		name string
		req  testutils.Req
		want testutils.Want
	}

	deps := dependencies.GetDependencies()
	tests := []testCase{
		{
			name: "returns 201 and short link for valid POST text/plain",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "POST",
				ContentType: "text/plain",
				Body:        "https://example.com",
			},
			want: testutils.Want{
				Status:      http.StatusCreated,
				Body:        deps.URLService.CreateShortLink("https://example.com"),
				ContentType: "text/plain",
			},
		},
		{
			name: "returns 400 for non-POST method",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "GET",
				ContentType: "text/plain",
				Body:        "https://example.com",
			},
			want: testutils.Want{
				Status: http.StatusMethodNotAllowed,
				Body:   "",
			},
		},
		{
			name: "returns 400 for wrong content type",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "POST",
				ContentType: "application/json",
				Body:        "https://example.com",
			},
			want: testutils.Want{
				Status: http.StatusBadRequest,
				Body:   "",
			},
		},
		{
			name: "returns 400 for empty body",
			req: testutils.Req{
				URL:         "/",
				MethodName:  "POST",
				ContentType: "text/plain",
				Body:        "",
			},
			want: testutils.Want{
				Status: http.StatusBadRequest,
				Body:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testutils.TestRequest(t, httptest.NewServer(deps.Router), tt.req)
			defer res.Body.Close()
			assert.Equal(t, tt.want.Status, res.StatusCode)
			assert.Equal(t, tt.want.ContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.Body, body)
		})
	}
}

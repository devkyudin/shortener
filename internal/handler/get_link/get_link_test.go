package get_link_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/dependencies"
	"github.com/devkyudin/shortener/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetLink(t *testing.T) {
	type testCase struct {
		name string
		req  testutils.Req
		want testutils.Want
	}

	deps := dependencies.GetDependencies()
	tests := []testCase{
		{
			name: "returns 307 and Location for existing id",
			req: testutils.Req{
				Url:         "/" + deps.URLService.CreateShortLink("https://example.com")[len("http://localhost:8080/"):],
				MethodName:  "GET",
				ContentType: "text/plain",
			},
			want: testutils.Want{
				Status:      http.StatusTemporaryRedirect,
				Location:    "https://example.com",
				ContentType: "text/plain",
			},
		},
		{
			name: "returns 400 for non-GET method",
			req: testutils.Req{
				Url:         "/" + deps.URLService.CreateShortLink("https://example.com")[len("http://localhost:8080/"):],
				MethodName:  "POST",
				ContentType: "text/plain",
			},
			want: testutils.Want{
				Status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "returns 400 for missing id",
			req: testutils.Req{
				Url:         "/",
				MethodName:  "GET",
				ContentType: "text/plain",
			},
			want: testutils.Want{
				Status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "returns 400 for unknown id",
			req: testutils.Req{
				Url:         "/unknown-id",
				MethodName:  "GET",
				ContentType: "text/plain",
			},
			want: testutils.Want{
				Status: http.StatusBadRequest,
			},
		},
	}

	ts := httptest.NewServer(deps.Router)
	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testutils.TestRequest(t, httptest.NewServer(deps.Router), tt.req)
			defer res.Body.Close()
			assert.Equal(t, tt.want.Status, res.StatusCode)
			assert.Equal(t, tt.want.ContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.Location, res.Header.Get("Location"))
			assert.Equal(t, tt.want.Body, body)
		})
	}
}

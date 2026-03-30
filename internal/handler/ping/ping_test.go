package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devkyudin/shortener/internal/dependencies"
	"github.com/devkyudin/shortener/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	type testCase struct {
		name string
		req  testutils.Req
		want testutils.Want
	}

	url := "/ping"
	deps := dependencies.GetDependencies()
	tests := []testCase{
		{
			name: "",
			req: testutils.Req{
				URL:        url,
				MethodName: "GET",
			},
			want: testutils.Want{
				Status: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := testutils.TestRequest(t, httptest.NewServer(*deps.Router), tt.req)
			defer res.Body.Close()
			assert.Equal(t, tt.want.Status, res.StatusCode)
			assert.Equal(t, tt.want.ContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.Body, body)
		})
	}
}

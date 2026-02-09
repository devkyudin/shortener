package testutils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devkyudin/shortener/internal/service"
	"github.com/stretchr/testify/require"
)

type Req struct {
	URL         string
	MethodName  string
	ContentType string
	Body        string
}

type Want struct {
	Status      int
	Body        string
	ContentType string
	Location    string
}

func TestRequest(t *testing.T, ts *httptest.Server, req Req) (*http.Response, string) {
	reader := io.NopCloser(io.Reader(strings.NewReader(req.Body)))
	r, err := http.NewRequest(req.MethodName, ts.URL+req.URL, reader)
	require.NoError(t, err)

	r.Header.Set("Content-Type", req.ContentType)
	res, err := ts.Client().Do(r)
	require.NoError(t, err)
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res, string(resp)
}

func CreateShortLinkSafe(s *service.URLService, originalURL string) string {
	result, err := s.CreateShortLink(originalURL)
	if err != nil {
		panic("Failed to create short link: " + err.Error())
	}

	return result
}

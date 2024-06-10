package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/stretchr/testify/assert"
)

var (
	testHeader = http.CanonicalHeaderKey("x-forwarded-for")
	realIp     = "192.168.19.10"
	proxyIp    = "10.0.0.2"
)

func TestProxyHandler(t *testing.T) {
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.RemoteAddr, realIp)
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = proxyIp
	req.Header = map[string][]string{
		testHeader: {realIp},
	}

	middleware := server.ProxyMiddleware(validationHandler)
	middleware.ServeHTTP(httptest.NewRecorder(), req)
}

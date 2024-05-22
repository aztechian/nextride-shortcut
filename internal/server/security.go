package server

import (
	"net/http"

	"github.com/unrolled/secure"
)

func SecurityMiddleware(next http.Handler) http.Handler {
	secure := secure.New(secure.Options{
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "script-src $NONCE",
	})
	return secure.Handler(next)
}

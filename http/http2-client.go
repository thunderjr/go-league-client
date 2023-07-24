package league_http

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"golang.org/x/net/http2"
)

type Response struct {
	StatusCode int `json:"status_code"`
}

func InitHttp2Client(certificate string) *http.Client {
	client := &http.Client{}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(certificate))

	client.Transport = &http2.Transport{TLSClientConfig: &tls.Config{
		RootCAs: caCertPool,
	}}

	return client
}

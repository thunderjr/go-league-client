package league_http

import (
	"crypto/tls"
	"net/http"

	league_auth "github.com/thunderjr/go-league-client/auth"
)

type LeagueHttp struct {
	credentials *league_auth.Credentials
	client      *http.Client
}

type LeagueClientOptions struct {
	Credentials *league_auth.Credentials
	UseHttp2    bool
}

func Init(options LeagueClientOptions) *LeagueHttp {
	if options.Credentials == nil {
		panic("credentials are required")
	}

	var client *http.Client

	if options.UseHttp2 {
		client = InitHttp2Client(options.Credentials.Certificate)
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	return &LeagueHttp{
		credentials: options.Credentials,
		client:      client,
	}
}

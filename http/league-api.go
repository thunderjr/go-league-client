package league_http

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
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

func (l *LeagueHttp) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://127.0.0.1:"+l.credentials.Port+path, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("riot:"+l.credentials.Password)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading response body: %s", err)
		return nil, err
	}

	return body, nil
}

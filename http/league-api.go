package league_http

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/thunderjr/go-league-client/auth"
)

type LeagueClient struct {
	credentials *auth.Credentials
	client      *http.Client
}

func NewLeagueClient(credentials *auth.Credentials) *LeagueClient {
	return &LeagueClient{
		client: InitHttp2Client(credentials.Certificate),
		// client: &http.Client{
		// 	Transport: &http.Transport{
		// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// 	},
		// },
		credentials: credentials,
	}
}

func (l *LeagueClient) Get(path string) ([]byte, error) {
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

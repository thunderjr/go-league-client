package league_http

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func (l *LeagueHttp) makeRequest(method, path string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, "https://127.0.0.1:"+l.credentials.Port+path, body)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, 0, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("riot:"+l.credentials.Password)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return nil, resp.StatusCode, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading response body: %s", err)
		return nil, resp.StatusCode, err
	}

	return responseBody, resp.StatusCode, nil
}

func (l *LeagueHttp) Get(path string) ([]byte, int, error) {
	return l.makeRequest("GET", path, nil)
}

func (l *LeagueHttp) Post(path string, body io.Reader) ([]byte, int, error) {
	return l.makeRequest("POST", path, body)
}

func (l *LeagueHttp) Patch(path string, body io.Reader) ([]byte, int, error) {
	return l.makeRequest("PATCH", path, body)
}

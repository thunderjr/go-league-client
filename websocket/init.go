package league_websocket

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	league_auth "github.com/thunderjr/go-league-client/auth"
)

type EventResponse struct {
	URI  string
	Data interface{}
}

type EventCallback func(event EventResponse)

type ConnectionOptions struct {
	Credentials  *league_auth.Credentials
	PollInterval time.Duration
	MaxRetries   int
}

type LeagueWebSocket struct {
	Conn          *websocket.Conn
	Subscriptions map[string][]EventCallback
}

func createWebSocketConnection(address string, authHeader http.Header, dialer *websocket.Dialer) (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	conn, _, err := dialer.DialContext(ctx, address, authHeader)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (lws *LeagueWebSocket) subscribeToJsonEvents() error {
	subscriptionMessage, err := json.Marshal([]interface{}{5, "OnJsonApiEvent"})
	if err != nil {
		return err
	}

	if err := lws.Conn.WriteMessage(websocket.TextMessage, subscriptionMessage); err != nil {
		return err
	}

	return nil
}

func Init(options *ConnectionOptions) (*LeagueWebSocket, error) {
	if options.Credentials == nil {
		return nil, errors.New("credentials are required")
	}

	var internalRetryCount = 0

	var tlsConfig *tls.Config
	if options.Credentials.Certificate != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(options.Credentials.Certificate))

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	} else {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	u := url.URL{Scheme: "wss", Host: fmt.Sprintf("127.0.0.1:%s", options.Credentials.Port)}

	authHeader := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte("riot:"+options.Credentials.Password))},
	}

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  tlsConfig,
	}

	conn, err := createWebSocketConnection(u.String(), authHeader, dialer)
	if err != nil {
		log.Println(err)
		if options.MaxRetries == 0 || (options.MaxRetries > 0 && internalRetryCount > options.MaxRetries) {
			return nil, errors.New("could not connect to lcu websocket api")
		}

		time.Sleep(options.PollInterval)
		internalRetryCount++

		return Init(options)
	}

	lws := &LeagueWebSocket{
		Conn:          conn,
		Subscriptions: make(map[string][]EventCallback),
	}

	if err := lws.subscribeToJsonEvents(); err != nil {
		return nil, errors.New("could not subscribe to json events")
	}

	go lws.listen()

	return lws, nil
}

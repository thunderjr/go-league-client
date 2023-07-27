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
	"reflect"
	"strings"
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

func (lws *LeagueWebSocket) listen() {
	for {
		_, message, err := lws.Conn.ReadMessage()
		if err != nil {
			log.Println("[listen] read ws message:", err)
			return
		}

		/*
		 * The websocket message response is a slice with the following (in order):
		 * 	- EventType int
		 * 	- EventName string (here it's always == 'OnJsonApiEvent')
		 * 	- Data map[string]interface{} (actual useful data)
		 *		= The keys for the Data map are:
		 *			- data
		 *			- eventType
		 *			- uri
		 */

		// res := make([]interface{}, 0, 3)
		res := []interface{}{
			0,                            // EventType
			"",                           // EventName
			make(map[string]interface{}), // Data
		}

		// The subscription to the OnJsonApiEvent returns an empty slice as response
		// taking care of it here
		if reflect.ValueOf(message).Kind() == reflect.Slice {
			slice := reflect.ValueOf(message)
			if slice.Len() == 0 {
				continue
			}
		}

		if err := json.Unmarshal(message, &res); err != nil {
			log.Println("[listen] unmarshal event:", err)
			continue
		}

		// Accessing the event data structure
		resData := res[2].(map[string]interface{})

		eventResponse := EventResponse{
			Data: resData["data"],
			URI:  resData["uri"].(string),
		}

		if callbacks, ok := lws.Subscriptions[eventResponse.URI]; ok {
			for _, cb := range callbacks {
				cb(eventResponse)
			}
		}
	}
}

func (lws *LeagueWebSocket) Subscribe(path string, effect EventCallback) {
	p := strings.Trim(path, " ")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	if _, ok := lws.Subscriptions[p]; !ok {
		lws.Subscriptions[p] = []EventCallback{effect}
	} else {
		lws.Subscriptions[p] = append(lws.Subscriptions[p], effect)
	}
}

func (lws *LeagueWebSocket) Unsubscribe(path string) {
	p := strings.Trim(path, " ")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	delete(lws.Subscriptions, p)
}

func Init(options *ConnectionOptions) (*LeagueWebSocket, error) {
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

	subscriptionMessage, err := json.Marshal([]interface{}{5, "OnJsonApiEvent"})
	if err != nil {
		return nil, err
	}

	if err := lws.Conn.WriteMessage(websocket.TextMessage, subscriptionMessage); err != nil {
		return nil, err
	}

	go lws.listen()

	return lws, nil
}

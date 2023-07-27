# Go League of Legends Client SDK

Inspired by [league-connect](https://github.com/matsjla/league-connect)

## Getting Started

[Official League Client API Documentation](https://developer.riotgames.com/docs/lol#league-client-api)

[HTTP and HTTP2 (HextechDocs)](https://hextechdocs.dev/getting-started-with-the-lcu-api/)

[WebSocket (HextechDocs)](https://hextechdocs.dev/getting-started-with-the-lcu-websocket/)

[League Client API FAQ](https://hextechdocs.dev/lcu-api-faq/)

## Installation

To install the package, run the following command:

```bash
go get github.com/thunderjr/go-league-client
```

## Usage

### WebSocket

```go
package main

import (
	"context"
	"fmt"
	"time"

	league_auth "github.com/thunderjr/go-league-client/auth"
	league_websocket "github.com/thunderjr/go-league-client/websocket"
)

func main() {
	ctx := context.Background()

	// Initialize the league auth struct
	auth := league_auth.Init(league_auth.AuthenticationOptions{
		AwaitConnection: true,
		Timeout:         time.Minute,
	})

	// Get the current local server credentials from leagueClient proccess
	auth.Authenticate(ctx)

	// Define your connection options
	options := &league_websocket.ConnectionOptions{
		// Use the LeagueClient credentials to connect to the WebSocket server
		Credentials:  auth.Credentials,
		PollInterval: time.Second,
		MaxRetries:   10,
	}

	// Initialize a WebSocket connection
	lws, err := league_websocket.Init(options)
	if err != nil {
		fmt.Println("Error initializing WebSocket connection:", err)
		return
	}

	// Define a callback function for the subscription
	callback := func(event league_websocket.EventResponse) {
		fmt.Println("Received data:", event.Data)
		fmt.Println("Event URI:", event.URI)
	}

	// Subscribe to the match updates URI
	// Create a room to receive some example data
	lws.Subscribe("/lol-gameflow/v1/session", callback)

	// Wait for a while to receive some data
	time.Sleep(60 * time.Second)

	// Unsubscribe
	lws.Unsubscribe("/lol-gameflow/v1/session")
}
```

Replace `"/lol-gameflow/v1/session"` with the actual path you want to subscribe to.

### HTTP

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	league_auth "github.com/thunderjr/go-league-client/auth"
	league_http "github.com/thunderjr/go-league-client/http"
)

func main() {
	ctx := context.Background()

	// Initialize and authenticate the League of Legends client
	auth := league_auth.Init(league_auth.AuthenticationOptions{
		AwaitConnection: true,
		Timeout:         time.Minute,
	})

	auth.Authenticate(ctx)

	// Initialize the League of Legends HTTP client
	httpClient := league_http.Init(league_http.LeagueClientOptions{
		Credentials: auth.Credentials,
		UseHttp2:    true,
	})

	// Make a GET request to retrieve the current summoner data
	endpoint := "/lol-summoner/v1/current-summoner"
	body, err := httpClient.Get(endpoint)
	if err != nil {
		log.Fatalln("Failed to make GET request:", err)
	}

	fmt.Println("Current Summoner Data:")
	fmt.Println(string(body))
}
```

Replace `"/lol-summoner/v1/current-summoner"` with the actual endpoint you want to request.

## [Examples](https://github.com/thunderjr/go-league-client/tree/master/examples)

## Contributing

Contributions are welcome! If you have a feature request, bug report, or proposal for code refactoring, please feel free to open an issue or submit a pull request.

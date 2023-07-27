package main

import (
	"context"
	"fmt"
	"time"

	league_auth "github.com/thunderjr/go-league-client/auth"
	league_websocket "github.com/thunderjr/go-league-client/websocket"
)

func LeagueWebSocketExample() {
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

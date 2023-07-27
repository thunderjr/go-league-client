package main

import (
	"context"
	"fmt"
	"log"
	"time"

	league_auth "github.com/thunderjr/go-league-client/auth"
	league_http "github.com/thunderjr/go-league-client/http"
)

func LeagueHttpExample() {
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

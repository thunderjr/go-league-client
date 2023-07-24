package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thunderjr/go-league-client/auth"
	league_http "github.com/thunderjr/go-league-client/http"
)

func main() {
	ctx := context.Background()

	leagueAuth := auth.LeagueAuth{
		AuthenticationOptions: auth.AuthenticationOptions{
			AwaitConnection: true,
			Timeout:         time.Minute,
		},
	}

	leagueAuth.Authenticate(ctx)

	httpClient := league_http.NewLeagueClient(leagueAuth.Credentials)

	endpoint := "/lol-summoner/v1/current-summoner"

	body, err := httpClient.Get(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))

	body, err = httpClient.Get(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
}

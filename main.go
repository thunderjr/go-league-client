package main

import (
	"context"
	"fmt"
	"time"

	"github.com/thunderjr/go-league-client/auth"
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

	fmt.Println(leagueAuth.Credentials)
}

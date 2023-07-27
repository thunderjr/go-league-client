package league_websocket_test

import (
	"testing"

	league_websocket "github.com/thunderjr/go-league-client/websocket"
)

func TestSubscribe(t *testing.T) {
	lws := &league_websocket.LeagueWebSocket{
		Subscriptions: make(map[string][]league_websocket.EventCallback),
	}

	var state int
	callback := func(response league_websocket.EventResponse) { state++ }

	lws.Subscribe("/test/path", callback)

	// Trigger the callback to change the state
	lws.Subscriptions["/test/path"][0](league_websocket.EventResponse{})
	if state != 1 {
		t.Error("Failed to subscribe to the new path")
	}

	// Subscribe to an existing path
	lws.Subscribe("/test/path", callback)

	// Trigger the callbacks to change the state
	for _, callback := range lws.Subscriptions["/test/path"] {
		callback(league_websocket.EventResponse{})
	}
	if state != 3 {
		t.Error("Failed to subscribe to an existing path")
	}
}

func TestUnsubscribe(t *testing.T) {
	lws := &league_websocket.LeagueWebSocket{
		Subscriptions: make(map[string][]league_websocket.EventCallback),
	}

	lws.Subscriptions["/test/path"] = []league_websocket.EventCallback{func(response league_websocket.EventResponse) {}}

	lws.Unsubscribe("/test/path")

	if _, ok := lws.Subscriptions["/test/path"]; ok {
		t.Error("Failed to unsubscribe")
	}
}

func TestUnsubscribeNonExistingPath(t *testing.T) {
	lws := &league_websocket.LeagueWebSocket{
		Subscriptions: make(map[string][]league_websocket.EventCallback),
	}

	lws.Unsubscribe("/non/existing/path")

	// Check if there's no change for a non-existing path
	if _, ok := lws.Subscriptions["/non/existing/path"]; ok {
		t.Error("Unsubscribed from a non-existing path")
	}
}

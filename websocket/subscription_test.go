package league_websocket

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	lws := &LeagueWebSocket{
		Subscriptions: make(map[string][]EventCallback),
	}

	var state int
	callback := func(response EventResponse) { state++ }

	lws.Subscribe("/test/path", callback)

	// Trigger the callback to change the state
	lws.Subscriptions["/test/path"][0](EventResponse{})
	if state != 1 {
		t.Error("Failed to subscribe to the new path")
	}

	// Subscribe to an existing path
	lws.Subscribe("/test/path", callback)

	// Trigger the callbacks to change the state
	for _, callback := range lws.Subscriptions["/test/path"] {
		callback(EventResponse{})
	}
	if state != 3 {
		t.Error("Failed to subscribe to an existing path")
	}
}

func TestUnsubscribe(t *testing.T) {
	lws := &LeagueWebSocket{
		Subscriptions: make(map[string][]EventCallback),
	}

	lws.Subscriptions["/test/path"] = []EventCallback{func(response EventResponse) {}}

	lws.Unsubscribe("/test/path")

	if _, ok := lws.Subscriptions["/test/path"]; ok {
		t.Error("Failed to unsubscribe")
	}
}

func TestUnsubscribeNonExistingPath(t *testing.T) {
	lws := &LeagueWebSocket{
		Subscriptions: make(map[string][]EventCallback),
	}

	lws.Unsubscribe("/non/existing/path")

	// Check if there's no change for a non-existing path
	if _, ok := lws.Subscriptions["/non/existing/path"]; ok {
		t.Error("Unsubscribed from a non-existing path")
	}
}

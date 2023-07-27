package league_websocket

import "strings"

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

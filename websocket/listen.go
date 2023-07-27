package league_websocket

import (
	"log"
	"reflect"
)

func (lws *LeagueWebSocket) listen() {
	for {
		_, message, err := lws.Conn.ReadMessage()
		if err != nil {
			log.Println("[listen] read ws message:", err)
			return
		}

		// The subscription to the OnJsonApiEvent returns an empty slice as response (taking care of it here)
		if reflect.ValueOf(message).Kind() == reflect.Slice {
			slice := reflect.ValueOf(message)
			if slice.Len() == 0 {
				continue
			}
		}

		eventResponse, err := parseMessage(message)
		if err != nil {
			log.Println("[listen] parse event:", err)
			continue
		}

		if callbacks, ok := lws.Subscriptions[eventResponse.URI]; ok {
			for _, cb := range callbacks {
				cb(eventResponse)
			}
		}
	}
}

package league_websocket

import (
	"encoding/json"
)

func ParseMessage(message []byte) (EventResponse, error) {
	/*
	 * The websocket message response is a slice with the following (in order):
	 * 	Ex.: [5, "OnJsonApiEvent", {"uri": "/some/uri", "data": {"key": "value"}}]
	 * 	- EventType int
	 * 	- EventName string (here it's always == 'OnJsonApiEvent')
	 * 	- Data map[string]interface{} (actual useful data)
	 *			= The keys for the Data map are:
	 *				- data
	 *				- eventType
	 *				- uri
	 */
	res := []interface{}{
		0,                            // EventType
		"",                           // EventName
		make(map[string]interface{}), // Data
	}

	if err := json.Unmarshal(message, &res); err != nil {
		return EventResponse{}, err
	}

	resData := res[2].(map[string]interface{})

	eventResponse := EventResponse{
		Data: resData["data"],
		URI:  resData["uri"].(string),
	}

	return eventResponse, nil
}

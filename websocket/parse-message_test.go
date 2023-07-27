package league_websocket

import (
	"testing"
)

func TestParseMessage(t *testing.T) {
	expectedEvent := EventResponse{
		URI: "/some/uri",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	result, err := parseMessage([]byte(`[5, "OnJsonApiEvent", {"uri": "/some/uri", "data": {"key": "value"}}]`))

	if err != nil {
		t.Errorf("parseMessage() error = %v", err)
		return
	}

	if result.URI != expectedEvent.URI {
		t.Errorf("parseMessage() got = %v, expected %v", result.URI, expectedEvent.URI)
	}

	gotData, ok1 := result.Data.(map[string]interface{})
	expectedData, ok2 := expectedEvent.Data.(map[string]interface{})

	if ok1 && ok2 {
		for key, expectedValue := range expectedData {
			if gotValue, ok := gotData[key]; !ok || gotValue != expectedValue {
				t.Errorf("parseMessage() got data = %v, expected data %v", result.Data, expectedEvent.Data)
			}
		}
	} else {
		t.Errorf("parseMessage() got data = %v, expected data %v", result.Data, expectedEvent.Data)
	}
}

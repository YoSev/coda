package utils

import (
	"encoding/json"
	"fmt"
)

func HandleJSON[T any](j json.RawMessage, handler func(*T) (json.RawMessage, error)) (json.RawMessage, error) {
	var params T
	if err := json.Unmarshal(j, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters, invalid parameters: %v", err)
	}
	return handler(&params)
}

func ReturnRaw(obj any) json.RawMessage {
	b, _ := json.Marshal(obj)
	return json.RawMessage([]byte(b))
}

func AddQuotes(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

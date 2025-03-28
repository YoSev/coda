package fn

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Fn struct{}

func New() *Fn {
	return &Fn{}
}

func handleJSON[T any](j json.RawMessage, handler func(*T) (json.RawMessage, error)) (json.RawMessage, error) {
	var params T
	if err := json.Unmarshal(j, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters, invalid parameters: %v", err)
	}
	return handler(&params)
}

func returnString(input string) json.RawMessage {
	return json.RawMessage(addQuotes(input))
}

func returnArray(input []string) json.RawMessage {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil
	}

	return json.RawMessage(jsonBytes)
}

func returnInt64(num int64) json.RawMessage {
	return json.RawMessage([]byte(strconv.FormatInt(num, 10)))
}

func returnAny(obj any) json.RawMessage {
	b, _ := json.Marshal(obj)
	return json.RawMessage([]byte(b))
}

func addQuotes(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

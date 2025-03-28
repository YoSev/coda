package utils

import "encoding/json"

// jsonEscape escapes the string to avoid errors when unmarshaling JSON later
func JsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}

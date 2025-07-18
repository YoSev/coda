package fn

import "testing"

func TestShoutrrr(t *testing.T) {
	var f *Fn
	var err error
	var result []byte

	f = New("0.0.0")
	fn := &fnMessage{category: FnCategoryMessaging}
	fn.init(f)

	// Example JSON input for Shoutrrr function
	jsonInput := `{"urls": ["telegram://123456789:TESTTOKEN987@telegram?chats=416898072"], "message": "Hello, *Shoutrrr*!", "parameters": {"parsemode": "markdown"}}`

	result, err = fn.shoutrrr([]byte(jsonInput))
	if err == nil {
		t.Fatalf("expected Unauthorized error: %v", err)
	}

	if len(result) > 0 {
		t.Fatalf("expected no result, got: %s", string(result))
	}
}

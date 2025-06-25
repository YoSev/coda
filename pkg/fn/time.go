package fn

import (
	"encoding/json"
	"time"
)

type FnTime struct {
	FnHandler
	*Fn

	Category FnCategory
}

func (f *FnTime) Register() {
	f.fns["time.datetime"] = &FnEntry{
		Handler:     f.generateDatetime,
		Name:        "Generate Datetime",
		Description: "Generates a datetime string based on the provided format",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The format for the datetime (e.g., '2006-01-02 15:04:05')", Mandatory: true},
		},
	}

	f.fns["time.timestamp.sec"] = &FnEntry{
		Handler:     f.generateTimestampSec,
		Name:        "Generate Timestamp (seconds)",
		Description: "Generates a timestamp in seconds",
		Category:    f.Category,
	}

	f.fns["time.timestamp.milli"] = &FnEntry{
		Handler:     f.generateTimestampMilli,
		Name:        "Generate Timestamp (milliseconds)",
		Description: "Generates a timestamp in milliseconds",
		Category:    f.Category,
	}

	f.fns["time.timestamp.micro"] = &FnEntry{
		Handler:     f.generateTimestampMicro,
		Name:        "Generate Timestamp (microseconds)",
		Description: "Generates a timestamp in microseconds",
		Category:    f.Category,
	}

	f.fns["time.timestamp.nano"] = &FnEntry{
		Handler:     f.generateTimestampNano,
		Name:        "Generate Timestamp (nanoseconds)",
		Description: "Generates a timestamp in nanoseconds",
		Category:    f.Category,
	}

	f.fns["time.sleep"] = &FnEntry{
		Handler:     f.sleep,
		Name:        "Sleep",
		Description: "Pauses execution for a specified duration in milliseconds",
		Category:    f.Category,
	}
}

type generateDatetimeParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *FnTime) generateDatetime(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *generateDatetimeParams) (json.RawMessage, error) {
		t := time.Now().Format(params.Value)
		return returnRaw(t), nil
	})
}
func (f *FnTime) generateTimestampSec(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().Unix()), nil
}
func (f *FnTime) generateTimestampMilli(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().UnixMilli()), nil
}
func (f *FnTime) generateTimestampMicro(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().UnixMicro()), nil
}
func (f *FnTime) generateTimestampNano(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().UnixNano()), nil
}

type sleepParams struct {
	Value int64 `json:"value" yaml:"value"`
}

func (f *FnTime) sleep(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *sleepParams) (json.RawMessage, error) {
		time.Sleep(time.Duration(params.Value) * time.Millisecond)
		return nil, nil
	})
}

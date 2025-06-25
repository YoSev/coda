package fn

import (
	"encoding/json"
	"time"

	"github.com/yosev/coda/internal/utils"
)

type fnTime struct {
	category FnCategory
}

func (f *fnTime) init(fn *Fn) {
	fn.register("time.datetime", &FnEntry{
		Handler:     f.generateDatetime,
		Name:        "Generate Datetime",
		Description: "Generates a datetime string based on the provided format",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The format for the datetime (e.g., '2006-01-02 15:04:05')", Mandatory: true},
		},
	})

	fn.register("time.timestamp.sec", &FnEntry{
		Handler:     f.generateTimestampSec,
		Name:        "Generate Timestamp (seconds)",
		Description: "Generates a timestamp in seconds",
		Category:    f.category,
	})

	fn.register("time.timestamp.milli", &FnEntry{
		Handler:     f.generateTimestampMilli,
		Name:        "Generate Timestamp (milliseconds)",
		Description: "Generates a timestamp in milliseconds",
		Category:    f.category,
	})

	fn.register("time.timestamp.micro", &FnEntry{
		Handler:     f.generateTimestampMicro,
		Name:        "Generate Timestamp (microseconds)",
		Description: "Generates a timestamp in microseconds",
		Category:    f.category,
	})

	fn.register("time.timestamp.nano", &FnEntry{
		Handler:     f.generateTimestampNano,
		Name:        "Generate Timestamp (nanoseconds)",
		Description: "Generates a timestamp in nanoseconds",
		Category:    f.category,
	})

	fn.register("time.sleep", &FnEntry{
		Handler:     f.sleep,
		Name:        "Sleep",
		Description: "Pauses execution for a specified duration in milliseconds",
		Category:    f.category,
	})
}

type generateDatetimeParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *fnTime) generateDatetime(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *generateDatetimeParams) (json.RawMessage, error) {
		t := time.Now().Format(params.Value)
		return utils.ReturnRaw(t), nil
	})
}
func (f *fnTime) generateTimestampSec(j json.RawMessage) (json.RawMessage, error) {
	return utils.ReturnRaw(time.Now().Unix()), nil
}
func (f *fnTime) generateTimestampMilli(j json.RawMessage) (json.RawMessage, error) {
	return utils.ReturnRaw(time.Now().UnixMilli()), nil
}
func (f *fnTime) generateTimestampMicro(j json.RawMessage) (json.RawMessage, error) {
	return utils.ReturnRaw(time.Now().UnixMicro()), nil
}
func (f *fnTime) generateTimestampNano(j json.RawMessage) (json.RawMessage, error) {
	return utils.ReturnRaw(time.Now().UnixNano()), nil
}

type sleepParams struct {
	Value int64 `json:"value" yaml:"value"`
}

func (f *fnTime) sleep(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *sleepParams) (json.RawMessage, error) {
		time.Sleep(time.Duration(params.Value) * time.Millisecond)
		return nil, nil
	})
}

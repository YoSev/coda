package fn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms/openai"
)

type openAI struct {
	ApiKey string `json:"api_key" yaml:"api_key"`
	Model  string `json:"model" yaml:"model"`
	Prompt string `json:"prompt" yaml:"prompt"`
	System string `json:"system" yaml:"system"`
}

func (f *Fn) OpenAI(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *openAI) (json.RawMessage, error) {
		apiKey := params.ApiKey
		modelName := params.Model

		llm, err := openai.New(
			openai.WithToken(apiKey),
			openai.WithModel(modelName),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize LLM: %v\n", err)
		}

		input := []openai.ChatMessage{}
		if params.System != "" {
			input = append(input, openai.ChatMessage{
				Role:    openai.RoleSystem,
				Content: params.System,
			})
		}
		input = append(input, openai.ChatMessage{
			Role:    openai.RoleUser,
			Content: params.Prompt,
		})

		// Serialize input to JSON or another string format
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize input: %v\n", err)
		}

		response, err := llm.Call(context.Background(), string(inputBytes))
		if err != nil {
			return nil, fmt.Errorf("error during LLM call: %v\n", err)
		}

		return returnRaw(response), nil
	})
}

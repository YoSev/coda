package fn

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type openAI struct {
	ApiKey      string   `json:"api_key" yaml:"api_key"`
	Model       string   `json:"model" yaml:"model"`
	Prompt      string   `json:"prompt" yaml:"prompt"`
	System      string   `json:"system" yaml:"system"`
	Attachments []string `json:"attachments" yaml:"attachments"`
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

		messages := []llms.MessageContent{}

		if params.System != "" {
			messages = append(messages, llms.MessageContent{
				Role:  llms.ChatMessageTypeSystem,
				Parts: []llms.ContentPart{llms.TextPart(params.System)},
			})
		}

		userParts := []llms.ContentPart{
			llms.TextPart(params.Prompt),
		}

		for _, attachment := range params.Attachments {
			if len(attachment) > 4 && attachment[:4] == "http" {
				userParts = append(userParts, llms.ImageURLPart(attachment))
			} else {
				fileContent, err := os.ReadFile(attachment)
				if err != nil {
					return nil, fmt.Errorf("failed to read attachment: %v\n", err)
				}
				userParts = append(userParts, llms.BinaryPart("application/octet-stream", fileContent))
			}
		}

		messages = append(messages, llms.MessageContent{
			Role:  llms.ChatMessageTypeHuman,
			Parts: userParts,
		})

		payload := map[string]interface{}{
			"model":    modelName,
			"messages": messages,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize payload: %v\n", err)
		}

		response, err := llm.Call(context.Background(), string(payloadBytes))
		if err != nil {
			return nil, fmt.Errorf("error during LLM call: %v\n", err)
		}

		return returnRaw(response), nil
	})
}

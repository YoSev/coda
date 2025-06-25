package fn

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/yosev/coda/internal/utils"
)

type fnAi struct {
	category FnCategory
}

func (f *fnAi) init(fn *Fn) {
	fn.register("ai.openai", &FnEntry{
		Handler:     f.openAI,
		Name:        "OpenAI",
		Description: "Performs an AI request",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "prompt", Description: "The actual prompt", Mandatory: true},
			{Name: "model", Description: "The modal to use", Mandatory: true},
			{Name: "api_key", Description: "The key to use", Mandatory: true},
			{Name: "system", Description: "The system query", Mandatory: false},
			{Name: "attachments", Description: "The attachments to include", Type: "array", Mandatory: false},
		},
	})
}

type openAIStruct struct {
	ApiKey      string   `json:"api_key" yaml:"api_key"`
	Model       string   `json:"model" yaml:"model"`
	Prompt      string   `json:"prompt" yaml:"prompt"`
	System      string   `json:"system" yaml:"system"`
	Attachments []string `json:"attachments" yaml:"attachments"`
}

func (f *fnAi) openAI(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *openAIStruct) (json.RawMessage, error) {
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

		return utils.ReturnRaw(response), nil
	})
}

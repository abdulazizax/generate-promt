package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func ExecutePrompt(openaiClient *openai.Client, prompt string) (string, error) {
	resp, err := openaiClient.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("Error executing prompt: %v", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("No choices returned in the response")
}

package ai

import (
	"context"
	"fmt"
	"generate-promt-v1/api/models"

	"github.com/sashabaranov/go-openai"
)

func ExecutePrompt(openaiClient *openai.Client, newPrompt string, conversationHistory *models.ConversationHistory) (string, error) {
	newMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: newPrompt,
	}
	conversationHistory.History = append(conversationHistory.History, newMessage)

	resp, err := openaiClient.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: conversationHistory.History,
	})
	if err != nil {
		return "", fmt.Errorf("Error executing prompt: %v", err)
	}

	if len(resp.Choices) > 0 {
		assistantResponse := resp.Choices[0].Message.Content
		conversationHistory.History = append(conversationHistory.History, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantResponse,
		})
		return assistantResponse, nil
	}
	return "", fmt.Errorf("No choices returned in the response")
}

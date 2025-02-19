package handlers

import (
	"generate-promt-v1/pkg/logger"

	"github.com/sashabaranov/go-openai"
)

type Handler struct {
	OpenaiClient *openai.Client
	Logger       logger.Logger
}

func New(openaiClient *openai.Client, logger logger.Logger) *Handler {
	return &Handler{
		OpenaiClient: openaiClient,
		Logger:       logger,
	}
}

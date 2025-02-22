package main

import (
	"generate-promt-v1/api"
	"generate-promt-v1/config"
	"generate-promt-v1/pkg/logger"

	"github.com/sashabaranov/go-openai"
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.LogLevel, "generate-promt-v1")

	openaiClient := openai.NewClient(cfg.AiKey)

	server := api.New(openaiClient, logger)

	server.Run(":" + cfg.HttpPort)
}

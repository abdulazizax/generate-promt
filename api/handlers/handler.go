package handlers

import (
	"generate-promt-v1/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/sheets/v4"
	"log"
)

type Handler struct {
	OpenaiClient *openai.Client
	Logger       logger.Logger
	SheetService sheets.Service
}

func New(openaiClient *openai.Client, logger logger.Logger, service *sheets.Service) *Handler {
	if openaiClient == nil || service == nil {
		log.Print("open ai or google sheet service is nil")
		return nil
	}

	return &Handler{
		OpenaiClient: openaiClient,
		Logger:       logger,
		SheetService: *service,
	}
}

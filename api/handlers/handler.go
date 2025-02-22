package handlers

import (
	"generate-promt-v1/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/sheets/v4"
	"log"
)

type Handler struct {
	OpenaiClient *openai.Client
	Logger       logger.Logger
	SheetService *sheets.Service
	DocService   *docs.Service
}

func New(openaiClient *openai.Client, logger logger.Logger, service *sheets.Service, docService *docs.Service) *Handler {
	if openaiClient == nil || service == nil {
		log.Print("open ai or google sheet service is nil")
		return nil
	}

	return &Handler{
		OpenaiClient: openaiClient,
		Logger:       logger,
		SheetService: service,
		DocService:   docService,
	}
}

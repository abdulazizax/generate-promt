package handlers

import (
	"fmt"
	"generate-promt-v1/api/models"
	"generate-promt-v1/pkg/ai"
	"generate-promt-v1/pkg/helper"
	"net/http"
	"os"
	"strings"

	"generate-promt-v1/config"

	"github.com/gin-gonic/gin"
)

// ExecutePrompt godoc
// @Summary Execute AI Prompt
// @Description Processes the project input and executes the AI prompt
// @Tags prompt
// @Accept json
// @Produce json
// @Param body body models.ProjectInput true "Project Input"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /prompt/execute [post]
func (h *Handler) ExecutePrompt(ctx *gin.Context) {
	var (
		body models.ProjectInput
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	conversationHistory := models.ConversationHistory{
		ProjectInput: body,
	}

	for i := 1; i <= 6; i++ {
		conversationHistory.Id = i
		resp, err := h.processPrompt(&conversationHistory)
		if err != nil {
			h.ReturnError(ctx, config.ErrorInternalServer, "Error while executing prompt", http.StatusInternalServerError)
			return
		}
		switch conversationHistory.Id {
		case 1:
			helper.ExportFunctionalRequirementsToSheet()
		}
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{Message: "Success!"})
}

func (h *Handler) processPrompt(conversationHistory *models.ConversationHistory) (string, error) {
	newPromt, err := h.readFile(fmt.Sprintf("pkg/promt/prompt_%d.txt", conversationHistory.Id))
	if err != nil {
		return "", err
	}

	if conversationHistory.Id == 1 {
		newPromt = strings.Replace(newPromt, "[company_name]", conversationHistory.ProjectInput.CompanyName, -1)
		newPromt = strings.Replace(newPromt, "[project_summary]", conversationHistory.ProjectInput.ProjectSummary, -1)
		newPromt = strings.Replace(newPromt, "[competitors]", strings.Join(conversationHistory.ProjectInput.Competitors, ", "), -1)
		newPromt = strings.Replace(newPromt, "[client_goals]", strings.Join(conversationHistory.ProjectInput.ClientGoals, ", "), -1)
		newPromt = strings.Replace(newPromt, "[target_audience]", conversationHistory.ProjectInput.TargetAudience, -1)
		newPromt = strings.Replace(newPromt, "[key_integrations]", strings.Join(conversationHistory.ProjectInput.KeyIntegrations, ", "), -1)
		newPromt = strings.Replace(newPromt, "[constraints]", conversationHistory.ProjectInput.Constraints, -1)
	}

	res, err := ai.ExecutePrompt(h.OpenaiClient, newPromt, conversationHistory)
	if err != nil {
		return "", err
	}

	return res, err
}

func (h *Handler) readFile(path string) (string, error) {
	fmt.Println("Read this file => ", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

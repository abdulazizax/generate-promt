package handlers

import (
	"generate-promt-v1/api/models"
	"generate-promt-v1/pkg/ai"
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

	resp, err := h.processPrompt(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error while executing prompt", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile("response.txt", []byte(resp), 0644)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error while writing to file", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{Message: "Success!"})
}

func (h *Handler) processPrompt(body *models.ProjectInput) (string, error) {
	data, err := os.ReadFile("example.txt")
	if err != nil {
		return "", err
	}

	req := string(data)

	req = strings.Replace(req, "[company_name]", body.CompanyName, -1)
	req = strings.Replace(req, "[project_summary]", body.ProjectSummary, -1)
	req = strings.Replace(req, "[competitors]", strings.Join(body.Competitors, ", "), -1)
	req = strings.Replace(req, "[client_goals]", strings.Join(body.ClientGoals, ", "), -1)
	req = strings.Replace(req, "[target_audience]", body.TargetAudience, -1)
	req = strings.Replace(req, "[key_integrations]", strings.Join(body.KeyIntegrations, ", "), -1)
	req = strings.Replace(req, "[constraints]", body.Constraints, -1)

	res, err := ai.ExecutePrompt(h.OpenaiClient, req)
	if err != nil {
		return "", err
	}

	return res, err
}

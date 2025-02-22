package handlers

import (
	"encoding/json"
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
		body     models.ProjectInput
		sheetUrl string
		docUrl   string
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	spreadsheet, err := helper.CreateNewSpreadsheet(h.SheetService, body.CompanyName, []string{"functionalReqs", "tab2", "tab3"})
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, fmt.Sprintf("Failed to create a new spreadsheet: %v", err), http.StatusInternalServerError)
		return
	}
	doc, err := helper.CreateNewDoc(h.DocService, body.CompanyName, []string{"tab1", "tab2"})
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, fmt.Sprintf("Failed to create a new doc: %v", err), http.StatusInternalServerError)
		return
	}
	conversationHistory := models.ConversationHistory{
		ProjectInput: body,
	}

	for i := 1; i <= 6; i++ {
		conversationHistory.Id = i
		resp, err := h.processPrompt(&conversationHistory)
		if err != nil {
			h.ReturnError(ctx, config.ErrorInternalServer, fmt.Sprintf("Error while executing prompt %v", err), http.StatusInternalServerError)
			return
		}
		switch conversationHistory.Id {
		case 1:
			var projectData models.ProjectResponse
			if err := json.Unmarshal([]byte(resp), &projectData); err != nil {
				h.ReturnError(ctx, config.ErrorBadRequest, fmt.Sprintf("%v", err), http.StatusBadRequest)
				return
			}
			if err := helper.ExportFunctionalRequirementsToSheet(h.SheetService, spreadsheet.SpreadsheetId, "functionalReqs", &projectData); err != nil {
				h.ReturnError(ctx, config.ErrorInternalServer, "failed to writing response to excel sheets", http.StatusInternalServerError)
				return
			}
			sheetUrl = fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", spreadsheet.SpreadsheetId)

			if err := helper.ExportProjectDataToDoc(h.DocService, doc.DocumentId, "tab1", &projectData.ProjectBrief); err != nil {
				h.ReturnError(ctx, config.ErrorInternalServer, "failed to writing response to google doc", http.StatusInternalServerError)
				return
			}
			docUrl = fmt.Sprintf("https://docs.google.com/document/d/%s/edit", doc.DocumentId)
		case 2:
		case 3:
		}
		break
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Success!",
		"docUrl":   docUrl,
		"sheetUrl": sheetUrl,
	})
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

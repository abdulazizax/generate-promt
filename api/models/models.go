package models

import "github.com/sashabaranov/go-openai"

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ProjectInput struct {
	CompanyName     string   `json:"company_name"`
	ProjectSummary  string   `json:"project_summary"`
	Competitors     []string `json:"competitors"`
	ClientGoals     []string `json:"client_goals"`
	TargetAudience  string   `json:"target_audience"`
	KeyIntegrations []string `json:"key_integrations"`
	Constraints     string   `json:"constraints"`
}

type ConversationHistory struct {
	Id           int
	ProjectInput ProjectInput
	History      []openai.ChatCompletionMessage
}

type ProjectResponse struct {
	ProjectBrief struct {
		ProjectGoal       string `json:"project_goal"`
		PrimaryObjectives string `json:"primary_objectives"`
		ExpectedOutcomes  string `json:"expected_outcomes"`
		SuccessMetrics    string `json:"success_metrics"`
	} `json:"project_brief"`

	FunctionalRequirements []struct {
		Epic    string `json:"epic"`
		Stories []struct {
			Story string   `json:"story"`
			Tasks []string `json:"tasks"`
		} `json:"stories"`
	} `json:"functional_requirements"`
}

package models

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

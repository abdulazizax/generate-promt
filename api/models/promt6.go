package models

type TeamMember struct {
	Role          string  `json:"role"`
	Count         float64 `json:"count"`
	Months        float64 `json:"months"`
	MonthlySalary float64 `json:"monthlySalary"`
	Sum           float64 `json:"sum"`
}

type Module struct {
	ModuleName string  `json:"moduleName"`
	Hours      float64 `json:"hours"`
	HourlyRate float64 `json:"hourlyRate"`
	Cost       float64 `json:"cost"`
}

type FinancialPlan struct {
	PrepaymentPercent float64   `json:"prepaymentPercent"` // Например, 30
	Prepayment        float64   `json:"prepayment"`        // Сумма предоплаты
	MonthlyPayments   []float64 `json:"monthlyPayments"`   // Платежи по месяцам
	TotalProjectCost  float64   `json:"totalProjectCost"`
}

type ProjectEstimate struct {
	Team          []TeamMember  `json:"team"`
	Modules       []Module      `json:"modules"`
	FinancialPlan FinancialPlan `json:"financialPlan"`
}

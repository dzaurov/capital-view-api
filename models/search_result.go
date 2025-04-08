// models/search_result.go (или models/dto.go)
package models

// FinancialReportDetail объединяет финансовые отчеты для конкретного периода/файла
type FinancialReportDetail struct {
	FinancialStatementInfo *FinancialStatement `json:"financial_statement_info,omitempty"`
	IncomeStatement        *IncomeStatement    `json:"income_statement,omitempty"`
	BalanceSheet           *BalanceSheet       `json:"balance_sheet,omitempty"`
	CashFlowStatement      *CashFlowStatement  `json:"cash_flow_statement,omitempty"`
}

// CompanySearchResult объединяет всю информацию по найденной компании
type CompanySearchResult struct {
	RegisterInfo     *Registers              `json:"register_info"` // Сделаем основным, не опциональным
	Members          []Member                `json:"members,omitempty"`
	BeneficialOwners []BeneficialOwner       `json:"beneficial_owners,omitempty"`
	FinancialReports []FinancialReportDetail `json:"financial_reports,omitempty"`
}

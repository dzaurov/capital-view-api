package models

type FinancialStatement struct {
	ID     uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID *string `json:"file_id,omitempty"`
	// Индекс нужен для поиска по компании и для связи в Registers
	LegalEntityRegistrationNumber *string `gorm:"uniqueIndex:uq_fs_company_year"` // <-- Составной уникальный индекс
	Year                          *string `gorm:"uniqueIndex:uq_fs_company_year"` // <-- Составной уникальный индекс
	// ... остальные поля ...
	SourceSchema     *string `json:"source_schema,omitempty"`
	SourceType       *string `json:"source_type,omitempty"`
	YearStartedOn    *string `json:"year_started_on,omitempty"`
	YearEndedOn      *string `json:"year_ended_on,omitempty"`
	Employees        *string `json:"employees,omitempty"`
	RoundedToNearest *string `json:"rounded_to_nearest,omitempty"`
	Currency         *string `json:"currency,omitempty"`
	CreatedAt        *string `json:"created_at,omitempty"`
	// --- Связи (Один к одному/нулю с деталями) ---
	// Указываем, что у Statement есть детали, связанные по ID этого Statement
	// и колонке 'statement_id' в таблицах деталей.
	IncomeStatement   *IncomeStatement   `gorm:"foreignKey:StatementID;references:ID"`
	BalanceSheet      *BalanceSheet      `gorm:"foreignKey:StatementID;references:ID"`
	CashFlowStatement *CashFlowStatement `gorm:"foreignKey:StatementID;references:ID"`
}

// TableName() не нужен, GORM по умолчанию сделает "financial_statements"

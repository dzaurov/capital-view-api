package models

type FinancialStatement struct {
	ID                            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID                        *string `json:"file_id,omitempty"`
	LegalEntityRegistrationNumber *string `gorm:"uniqueIndex:uq_fs_company_year"` // <--- Добавить/обновить тег
	SourceSchema                  *string `json:"source_schema,omitempty"`
	SourceType                    *string `json:"source_type,omitempty"`
	Year                          *string `gorm:"uniqueIndex:uq_fs_company_year"` // <--- Добавить/обновить тег
	YearStartedOn                 *string `json:"year_started_on,omitempty"`
	YearEndedOn                   *string `json:"year_ended_on,omitempty"`
	Employees                     *string `json:"employees,omitempty"`
	RoundedToNearest              *string `json:"rounded_to_nearest,omitempty"`
	Currency                      *string `json:"currency,omitempty"`
	CreatedAt                     *string `json:"created_at,omitempty"`
}

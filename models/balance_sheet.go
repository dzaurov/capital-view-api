package models

type BalanceSheet struct {
	ID                           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	StatementID                  *string `gorm:"uniqueIndex"` // <--- Добавить/обновить тег
	FileID                       *string `json:"file_id,omitempty"`
	Cash                         *string `json:"cash,omitempty"`
	MarketableSecurities         *string `json:"marketable_securities,omitempty"`
	AccountsReceivable           *string `json:"accounts_receivable,omitempty"`
	Inventories                  *string `json:"inventories,omitempty"`
	TotalCurrentAssets           *string `json:"total_current_assets,omitempty"`
	Investments                  *string `json:"investments,omitempty"`
	FixedAssets                  *string `json:"fixed_assets,omitempty"`
	IntangibleAssets             *string `json:"intangible_assets,omitempty"`
	TotalNonCurrentAssets        *string `json:"total_non_current_assets,omitempty"`
	TotalAssets                  *string `json:"total_assets,omitempty"`
	FutureHousingRepairsPayments *string `json:"future_housing_repairs_payments,omitempty"`
	CurrentLiabilities           *string `json:"current_liabilities,omitempty"`
	NonCurrentLiabilities        *string `json:"non_current_liabilities,omitempty"`
	Provisions                   *string `json:"provisions,omitempty"`
	Equity                       *string `json:"equity,omitempty"`
	TotalEquities                *string `json:"total_equities,omitempty"`
}

// models/registers.go
package models

type Registers struct {
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	Regcode            *string `gorm:"uniqueIndex"` // Уникальный индекс уже есть
	Sepa               *string
	Name               *string `gorm:"index"` // <-- Добавим индекс для поиска по имени
	NameBeforeQuotes   *string `gorm:"column:name_before_quotes"`
	NameInQuotes       *string `gorm:"column:name_in_quotes;index"` // <-- Добавим индекс
	NameAfterQuotes    *string `gorm:"column:name_after_quotes"`
	WithoutQuotes      *string `gorm:"column:without_quotes;index"` // <-- Добавим индекс
	Regtype            *string
	RegtypeText        *string `gorm:"column:regtype_text"`
	Type               *string
	TypeText           *string `gorm:"column:type_text"`
	Registered         *string
	Terminated         *string
	Closed             *string
	Address            *string
	IndexCompany       *string `gorm:"column:index_company"`
	Addressid          *string
	Region             *string
	City               *string
	Atvk               *string
	ReregistrationTerm *string `gorm:"column:reregistration_term"`
	Latitude           *string
	Longitude          *string

	// --- Связи ---
	// Указываем, что у одной записи Registers может быть много записей Member,
	// связанных по колонке 'legal_entity_registration_number' в таблице members,
	// которая соответствует колонке 'regcode' в таблице registers.
	Members             []Member             `gorm:"foreignKey:LegalEntityRegistrationNumber;references:Regcode"`
	BeneficialOwners    []BeneficialOwner    `gorm:"foreignKey:LegalEntityRegistrationNumber;references:Regcode"`
	FinancialStatements []FinancialStatement `gorm:"foreignKey:LegalEntityRegistrationNumber;references:Regcode"`
}

// Метод TableName оставляем, чтобы гарантировать имя "registers"
func (Registers) TableName() string {
	return "registers"
}

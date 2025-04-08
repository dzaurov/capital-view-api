package models

type BeneficialOwner struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// Индекс нужен для поиска по компании и для связи в Registers
	LegalEntityRegistrationNumber *string `gorm:"index" json:"legal_entity_registration_number,omitempty"` // <-- Индекс для связи
	Forename                      *string `gorm:"index" json:"forename,omitempty"`                         // <-- Индекс для поиска
	Surname                       *string `gorm:"index" json:"surname,omitempty"`                          // <-- Индекс для поиска
	// ... остальные поля ...
	LatvianIdentityNumberMasked *string `json:"latvian_identity_number_masked,omitempty"`
	BirthDate                   *string `json:"birth_date,omitempty"`
	Nationality                 *string `json:"nationality,omitempty"`
	Residence                   *string `json:"residence,omitempty"`
	RegisteredOn                *string `json:"registered_on,omitempty"`
	LastModifiedAt              *string `json:"last_modified_at,omitempty"`
	// Добавьте остальные поля из вашего CSV/модели, если они есть
}

// TableName() не нужен, GORM по умолчанию сделает "beneficial_owners"

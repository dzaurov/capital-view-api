package models

type Member struct {
	ID  uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Uri *string `json:"uri,omitempty"`
	// Индекс нужен для поиска по компании и для связи в Registers
	AtLegalEntityRegistrationNumber *string `gorm:"column:at_legal_entity_registration_number;index" json:"at_legal_entity_registration_number,omitempty"`
	EntityType                      *string `json:"entity_type,omitempty"`
	Name                            *string `gorm:"index" json:"name,omitempty"` // <-- Индекс для поиска по имени
	// Если это поле = regcode компании, где лицо является участником, то нужен индекс
	LegalEntityRegistrationNumber *string `gorm:"index" json:"legal_entity_registration_number,omitempty"` // <-- Индекс для связи
	// ... остальные поля ...
	LatvianIdentityNumberMasked *string `json:"latvian_identity_number_masked,omitempty"`
	BirthDate                   *string `json:"birth_date,omitempty"`
	NumberOfShares              *string `json:"number_of_shares,omitempty"`
	ShareNominalValue           *string `json:"share_nominal_value,omitempty"`
	ShareCurrency               *string `json:"share_currency,omitempty"`
	DateFrom                    *string `json:"date_from,omitempty"`
	RegisteredOn                *string `json:"registered_on,omitempty"`
	LastModifiedAt              *string `json:"last_modified_at,omitempty"`
	// Добавьте остальные поля из вашего CSV/модели, если они есть
}

// TableName() не нужен, GORM по умолчанию сделает "members"

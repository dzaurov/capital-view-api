package models

type Member struct {
	ID                              uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Uri                             *string `json:"uri,omitempty"`
	AtLegalEntityRegistrationNumber *string `json:"at_legal_entity_registration_number,omitempty"`
	EntityType                      *string `json:"entity_type,omitempty"`
	Name                            *string `json:"name,omitempty"`
	LatvianIdentityNumberMasked     *string `json:"latvian_identity_number_masked,omitempty"`
	BirthDate                       *string `json:"birth_date,omitempty"`
	LegalEntityRegistrationNumber   *string `json:"legal_entity_registration_number,omitempty"`
	NumberOfShares                  *string `json:"number_of_shares,omitempty"`
	ShareNominalValue               *string `json:"share_nominal_value,omitempty"`
	ShareCurrency                   *string `json:"share_currency,omitempty"`
	DateFrom                        *string `json:"date_from,omitempty"`
	RegisteredOn                    *string `json:"registered_on,omitempty"`
	LastModifiedAt                  *string `json:"last_modified_at,omitempty"`
}

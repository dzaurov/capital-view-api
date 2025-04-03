package models

type BeneficialOwner struct {
	ID                            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	LegalEntityRegistrationNumber *string `json:"legal_entity_registration_number,omitempty"`
	Forename                      *string `json:"forename,omitempty"`
	Surname                       *string `json:"surname,omitempty"`
	LatvianIdentityNumberMasked   *string `json:"latvian_identity_number_masked,omitempty"`
	BirthDate                     *string `json:"birth_date,omitempty"`
	Nationality                   *string `json:"nationality,omitempty"`
	Residence                     *string `json:"residence,omitempty"`
	RegisteredOn                  *string `json:"registered_on,omitempty"`
	LastModifiedAt                *string `json:"last_modified_at,omitempty"`
}

package models

type Register struct {
	ID                 uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Regcode            *string  `json:"regcode,omitempty"`
	Sepa               *string  `json:"sepa,omitempty"`
	Name               *string  `json:"name,omitempty"`
	NameBeforeQuotes   *string  `json:"name_before_quotes,omitempty"`
	NameInQuotes       *string  `json:"name_in_quotes,omitempty"`
	NameAfterQuotes    *string  `json:"name_after_quotes,omitempty"`
	WithoutQuotes      *string  `json:"without_quotes,omitempty"`
	Regtype            *string  `json:"regtype,omitempty"`
	RegtypeText        *string  `json:"regtype_text,omitempty"`
	Type               *string  `json:"type,omitempty"`
	TypeText           *string  `json:"type_text,omitempty"`
	Registered         *string  `json:"registered,omitempty"`
	Terminated         *string  `json:"terminated,omitempty"`
	Closed             *string  `json:"closed,omitempty"`
	Address            *string  `json:"address,omitempty"`
	IndexCompany       *string  `json:"index_company,omitempty"`
	Addressid          *string  `json:"addressid,omitempty"`
	Region             *string  `json:"region,omitempty"`
	City               *string  `json:"city,omitempty"`
	Atvk               *string  `json:"atvk,omitempty"`
	ReregistrationTerm *string  `json:"reregistration_term,omitempty"`
	Latitude           *float64 `json:"latitude,omitempty"`
	Longitude          *float64 `json:"longitude,omitempty"`
}

// models/registers.go (или где определена структура Registers)
package models

type Registers struct {
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	Regcode            *string `gorm:"uniqueIndex"`
	Sepa               *string
	Name               *string
	NameBeforeQuotes   *string `gorm:"column:name_before_quotes"`
	NameInQuotes       *string `gorm:"column:name_in_quotes"`
	NameAfterQuotes    *string `gorm:"column:name_after_quotes"`
	WithoutQuotes      *string `gorm:"column:without_quotes"`
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
}

// !!! --- ДОБАВЬТЕ ЭТОТ МЕТОД --- !!!
// Принудительно указывает GORM использовать таблицу "registers"
func (Registers) TableName() string {
	return "registers"
}

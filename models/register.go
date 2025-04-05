// models/register.go
package models

// Ваша СУЩЕСТВУЮЩАЯ структура Register
type Register struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	// Используем *string или string в зависимости от того, могут ли поля быть NULL
	// Убедитесь, что типы и теги gorm:"column:..." соответствуют ВАШЕЙ таблице register
	Regcode            *string `gorm:"uniqueIndex"` // Добавим uniqueIndex, как в DDL
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
	Index              *string `gorm:"column:index"` // Используем имя "index" из DDL таблицы register
	Addressid          *string
	Region             *string
	City               *string
	Atvk               *string
	ReregistrationTerm *string `gorm:"column:reregistration_term"`
	Latitude           *string // Оставляем TEXT, как в DDL таблицы register
	Longitude          *string // Оставляем TEXT, как в DDL таблицы register
}

// !!! --- ДОБАВЬТЕ ЭТОТ МЕТОД --- !!!
// Он явно укажет GORM использовать таблицу 'register' (в единственном числе)
// для структуры Register.
func (Register) TableName() string {
	return "register"
}

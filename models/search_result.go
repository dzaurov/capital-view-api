// models/search_result.go (ИЛИ models/models.go)
package models

// SimpleRegisterInfo содержит минимальный набор полей для списка результатов поиска
type SimpleRegisterInfo struct {
	Regcode     *string `json:"Regcode,omitempty"`
	Name        *string `json:"Name,omitempty"`
	RegtypeText *string `json:"RegtypeText,omitempty"`
	Address     *string `json:"Address,omitempty"`
	TypeText    *string `json:"TypeText,omitempty"`
}

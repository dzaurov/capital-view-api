// models/response.go
package models

// PaginatedResponse - стандартная структура для ответов с пагинацией
type PaginatedResponse struct {
	TotalRecords int64       `json:"total_records"` // Общее количество записей (не только на странице)
	Page         int         `json:"page"`          // Текущий номер страницы
	Limit        int         `json:"limit"`         // Лимит записей на странице
	Data         interface{} `json:"data"`          // Срез данных для текущей страницы
}

// HTTPError - структура для стандартной ошибки API (если еще не определена)
type HTTPError struct {
	Error string `json:"error"`
}

// NewHTTPError - конструктор для HTTPError (если еще не определен)
func NewHTTPError(err error) HTTPError {
	return HTTPError{Error: err.Error()}
}

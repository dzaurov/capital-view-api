// utils/pagination.go
package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20  // Количество записей на страницу по умолчанию
	MaxLimit     = 100 // Максимальное количество записей на страницу
)

// PaginationParams содержит разобранные параметры пагинации
type PaginationParams struct {
	Limit  int
	Offset int
	Page   int
}

// GetPaginationParams извлекает параметры 'page' и 'limit' из запроса,
// устанавливает значения по умолчанию и рассчитывает offset.
func GetPaginationParams(c *gin.Context) PaginationParams {
	pageStr := c.DefaultQuery("page", strconv.Itoa(DefaultPage))
	limitStr := c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DefaultPage
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = DefaultLimit
	}

	// Ограничиваем максимальный лимит
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
		Page:   page,
	}
}

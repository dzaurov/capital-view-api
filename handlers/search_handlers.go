// handlers/search_handlers.go
package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	// "sync" // sync.WaitGroup и sync.Mutex больше не нужны для Этапа 1

	"capital-view-api/db"     // <-- Убедитесь, что путь правильный
	"capital-view-api/models" // <-- Убедитесь, что путь правильный
	"capital-view-api/utils"  // <-- Убедитесь, что путь правильный

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DetailedSearch godoc
// @Summary Упрощенный поиск компаний (с пагинацией)
// @Description Ищет компании **только** по полям таблицы регистра (Regcode, SEPA, Name). Возвращает пагинированный список с базовой информацией. // <-- Описание изменено
// @Tags search
// @Produce json
// @Param q query string true "Поисковый запрос (Regcode, SEPA, Name)"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Записей на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.PaginatedResponse{data=[]models.SimpleRegisterInfo} "Пагинированный список базовой информации о компаниях"
// @Failure 400 {object} HTTPError "Неверный запрос (отсутствует 'q')"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /search/detailed [get]
func DetailedSearch(c *gin.Context) {
	log.Println("DetailedSearch (Simplified): Начало обработки запроса")

	// --- Параметры ---
	searchTerm := c.Query("q")
	if strings.TrimSpace(searchTerm) == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("поисковый параметр 'q' обязателен")))
		return
	}
	searchTermLower := strings.ToLower(searchTerm)
	searchTermLikeLower := "%" + searchTermLower + "%"
	pagination := utils.GetPaginationParams(c)
	log.Printf("DetailedSearch (Simplified): SearchTerm: '%s', Page: %d, Limit: %d", searchTerm, pagination.Page, pagination.Limit)

	// --- Этап 1 (Упрощен): Поиск ID ТОЛЬКО в таблице 'registers' ---
	log.Println("DetailedSearch (Simplified): Phase 1 starting - Finding matching IDs in 'registers' table only...")
	var registerMatches []struct{ Regcode *string } // Срез для получения только regcode

	// Выполняем ОДИН запрос к таблице registers
	err := db.DB.Model(&models.Registers{}). // Убедитесь, что используется модель для таблицы 'registers'
							Select("DISTINCT regcode"). // Выбираем УНИКАЛЬНЫЕ regcode
							Where("LOWER(regcode) = ?", searchTermLower).
							Or("LOWER(sepa) = ?", searchTermLower).
							Or("LOWER(name) LIKE ?", searchTermLikeLower).
							Or("LOWER(name_in_quotes) LIKE ?", searchTermLikeLower).
							Or("LOWER(without_quotes) LIKE ?", searchTermLikeLower).
							Find(&registerMatches).Error // Ищем все совпадающие regcode

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("DetailedSearch (Simplified): Error during ID search: %v", err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(fmt.Errorf("ошибка при поиске ID: %w", err)))
		return
	}

	// Собираем уникальные ID (хотя SELECT DISTINCT уже должен был это сделать)
	var uniqueRegCodes []string
	if registerMatches != nil { // Проверка на случай, если Find вернул nil при gorm.ErrRecordNotFound
		for _, match := range registerMatches {
			if match.Regcode != nil && *match.Regcode != "" {
				// Дополнительная проверка на уникальность не нужна из-за SELECT DISTINCT
				uniqueRegCodes = append(uniqueRegCodes, *match.Regcode)
			}
		}
	}
	sort.Strings(uniqueRegCodes) // Сортируем для стабильности пагинации
	totalRecords := int64(len(uniqueRegCodes))
	log.Printf("DetailedSearch (Simplified): Phase 1 finished - Found %d unique matching IDs.", totalRecords)

	// --- Пагинация найденных ID ---
	start := pagination.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(uniqueRegCodes) {
		start = len(uniqueRegCodes)
	}
	end := start + pagination.Limit
	if end > len(uniqueRegCodes) {
		end = len(uniqueRegCodes)
	}
	paginatedRegCodes := uniqueRegCodes[start:end]
	log.Printf("DetailedSearch (Simplified): Processing %d IDs for page %d", len(paginatedRegCodes), pagination.Page)

	// --- Создаем ответ по умолчанию ---
	simplePaginatedData := []models.SimpleRegisterInfo{} // Срез для УПРОЩЕННЫХ данных
	response := models.PaginatedResponse{
		TotalRecords: totalRecords, Page: pagination.Page, Limit: pagination.Limit, Data: simplePaginatedData,
	}

	if len(paginatedRegCodes) == 0 {
		log.Println("DetailedSearch (Simplified): No IDs to process for the current page or no results found.")
		c.JSON(http.StatusOK, response)
		return
	}

	// --- Этап 2: Загрузка ТОЛЬКО НЕОБХОДИМЫХ полей для пагинированных ID ---
	log.Printf("DetailedSearch (Simplified): Phase 2 starting - Fetching simplified data for %d IDs...", len(paginatedRegCodes))

	// Выбираем нужные поля в структуру SimpleRegisterInfo
	err = db.DB.Model(&models.Registers{}).
		Select("regcode", "name", "regtype_text", "address", "type_text"). // <--- ТОЛЬКО эти поля
		Where("regcode IN ?", paginatedRegCodes).
		// Order(...) // Можно добавить сортировку по regcode, чтобы соответствовать paginatedRegCodes
		Find(&simplePaginatedData).Error // Записываем результат в срез []SimpleRegisterInfo

	if err != nil {
		log.Printf("DetailedSearch (Simplified): Error fetching simplified data: %v", err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(fmt.Errorf("ошибка загрузки списка компаний: %w", err)))
		return
	}
	log.Printf("DetailedSearch (Simplified): Phase 2 finished - Fetched simplified data for %d companies.", len(simplePaginatedData))

	// --- Этап 3: Формирование финального ответа ---
	response.Data = simplePaginatedData // Обновляем поле Data

	log.Printf("DetailedSearch (Simplified): Request successful. Returning %d records.", len(simplePaginatedData))
	c.JSON(http.StatusOK, response)
}

// handlers/search_handlers.go
package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"capital-view-api/db"     // <-- Убедитесь, что путь правильный
	"capital-view-api/models" // <-- Убедитесь, что путь правильный
	"capital-view-api/utils"  // <-- Убедитесь, что путь правильный

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// -------------------------------------------------------------

// DetailedSearch godoc
// @Summary Расширенный поиск компаний (Оптимизированный + Пагинация + Вложенный ответ)
// @Description Ищет компании по части названия, Regcode, SEPA, имени бенефициара или участника. Возвращает пагинированный список companies (registers) с вложенными связанными данными.
// @Tags search
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Записей на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.PaginatedResponse{data=[]models.Registers} "Пагинированный список компаний (с вложенными Members, BeneficialOwners, FinancialStatements)" // <-- ИЗМЕНЕНО data type на models.Registers
// @Failure 400 {object} HTTPError "Неверный запрос (отсутствует 'q')"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /search/detailed [get]
func DetailedSearch(c *gin.Context) {
	log.Println("DetailedSearch: Начало обработки запроса")

	// --- Параметры ---
	searchTerm := c.Query("q")
	if strings.TrimSpace(searchTerm) == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("поисковый параметр 'q' обязателен")))
		return
	}
	searchTermLower := strings.ToLower(searchTerm)
	searchTermLikeLower := "%" + searchTermLower + "%"
	pagination := utils.GetPaginationParams(c)
	log.Printf("DetailedSearch: SearchTerm: '%s', Page: %d, Limit: %d", searchTerm, pagination.Page, pagination.Limit)

	// --- Этап 1: Поиск ВСЕХ совпадающих ID ---
	log.Println("DetailedSearch: Phase 1 starting - Finding all matching IDs...")
	var allFoundRegCodes []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	var searchErrors []error

	// Горутины поиска (registers, members, beneficial_owners)
	wg.Add(3)
	// Поиск в 'registers'
	go func() {
		defer wg.Done()
		var registerMatches []struct{ Regcode *string }
		err := db.DB.Model(&models.Registers{}).Select("regcode").
			Where("LOWER(regcode) = ?", searchTermLower).
			Or("LOWER(sepa) = ?", searchTermLower).
			Or("LOWER(name) LIKE ?", searchTermLikeLower).
			Or("LOWER(name_in_quotes) LIKE ?", searchTermLikeLower).
			Or("LOWER(without_quotes) LIKE ?", searchTermLikeLower).
			Find(&registerMatches).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := fmt.Errorf("register search error: %w", err)
			mu.Lock()
			searchErrors = append(searchErrors, errMsg)
			mu.Unlock()
			return
		}
		var localCodes []string
		for _, match := range registerMatches {
			if match.Regcode != nil && *match.Regcode != "" {
				localCodes = append(localCodes, *match.Regcode)
			}
		}
		if len(localCodes) > 0 {
			mu.Lock()
			allFoundRegCodes = append(allFoundRegCodes, localCodes...)
			mu.Unlock()
		}
	}()
	// Поиск в 'members'
	go func() {
		defer wg.Done()
		var memberMatches []struct{ LegalEntityRegistrationNumber *string }
		err := db.DB.Model(&models.Member{}).Select("legal_entity_registration_number").Where("LOWER(name) LIKE ?", searchTermLikeLower).Find(&memberMatches).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := fmt.Errorf("member search error: %w", err)
			mu.Lock()
			searchErrors = append(searchErrors, errMsg)
			mu.Unlock()
			return
		}
		var localCodes []string
		for _, match := range memberMatches {
			if match.LegalEntityRegistrationNumber != nil && *match.LegalEntityRegistrationNumber != "" {
				localCodes = append(localCodes, *match.LegalEntityRegistrationNumber)
			}
		}
		if len(localCodes) > 0 {
			mu.Lock()
			allFoundRegCodes = append(allFoundRegCodes, localCodes...)
			mu.Unlock()
		}
	}()
	// Поиск в 'beneficial_owners'
	go func() {
		defer wg.Done()
		var ownerMatches []struct{ LegalEntityRegistrationNumber *string }
		err := db.DB.Model(&models.BeneficialOwner{}).Select("legal_entity_registration_number").Where("LOWER(forename) LIKE ?", searchTermLikeLower).Or("LOWER(surname) LIKE ?", searchTermLikeLower).Find(&ownerMatches).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := fmt.Errorf("beneficial_owner search error: %w", err)
			mu.Lock()
			searchErrors = append(searchErrors, errMsg)
			mu.Unlock()
			return
		}
		var localCodes []string
		for _, match := range ownerMatches {
			if match.LegalEntityRegistrationNumber != nil && *match.LegalEntityRegistrationNumber != "" {
				localCodes = append(localCodes, *match.LegalEntityRegistrationNumber)
			}
		}
		if len(localCodes) > 0 {
			mu.Lock()
			allFoundRegCodes = append(allFoundRegCodes, localCodes...)
			mu.Unlock()
		}
	}()
	wg.Wait()
	log.Println("DetailedSearch: Phase 1 finished - All goroutines completed.")

	// Обработка ошибок поиска ID
	if len(searchErrors) > 0 {
		log.Printf("DetailedSearch: Error during ID search phase: %v", searchErrors)
		var errorStrings []string
		for _, e := range searchErrors {
			errorStrings = append(errorStrings, e.Error())
		}
		combinedError := fmt.Errorf("ошибки при поиске ID: %s", strings.Join(errorStrings, "; "))
		c.JSON(http.StatusInternalServerError, NewHTTPError(combinedError))
		return
	}

	// --- Уникализация и пагинация найденных ID ---
	uniqueRegCodesMap := make(map[string]struct{})
	var uniqueRegCodes []string
	for _, code := range allFoundRegCodes {
		if _, exists := uniqueRegCodesMap[code]; !exists {
			uniqueRegCodesMap[code] = struct{}{}
			uniqueRegCodes = append(uniqueRegCodes, code)
		}
	}
	sort.Strings(uniqueRegCodes)
	totalRecords := int64(len(uniqueRegCodes))
	log.Printf("DetailedSearch: Found %d unique matching IDs.", totalRecords)

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
	log.Printf("DetailedSearch: Processing %d IDs for page %d", len(paginatedRegCodes), pagination.Page)

	// --- Создаем пустой ответ по умолчанию ---
	// Мы будем возвращать []models.Registers напрямую
	paginatedData := []models.Registers{}
	response := models.PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
		Data:         paginatedData, // Начинаем с пустого среза нужного типа
	}

	// Если нет ID для обработки (на этой странице или вообще), возвращаем пустой ответ
	if len(paginatedRegCodes) == 0 {
		log.Println("DetailedSearch: No IDs to process for the current page.")
		c.JSON(http.StatusOK, response) // Возвращаем структуру с пустым Data
		return
	}

	// --- Этап 2 (Preload): Загрузка данных с предзагрузкой ---
	log.Printf("DetailedSearch: Phase 2 starting - Preloading data for %d IDs...", len(paginatedRegCodes))
	// Результат Preload будет записан в paginatedData (это []models.Registers)
	err := db.DB.
		Preload("Members").
		Preload("BeneficialOwners").
		Preload("FinancialStatements", func(db *gorm.DB) *gorm.DB { return db.Order("financial_statements.year DESC") }).
		Preload("FinancialStatements.IncomeStatement").
		Preload("FinancialStatements.BalanceSheet").
		Preload("FinancialStatements.CashFlowStatement").
		Where("regcode IN ?", paginatedRegCodes).
		// Попробуем добавить сортировку, чтобы соответствовать порядку paginatedRegCodes
		// Это может быть неэффективно для SQLite, но попробуем
		// Order(clause.OrderByColumn{Column: clause.Column{Name: "regcode"}, Values: paginatedRegCodes}). // Эта конструкция может не работать с IN, GORM сложен
		Find(&paginatedData).Error // <-- ЗАПИСЫВАЕМ РЕЗУЛЬТАТ НАПРЯМУЮ В paginatedData

	if err != nil {
		log.Printf("DetailedSearch: Error during Preload/Find: %v", err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(fmt.Errorf("ошибка загрузки данных компаний: %w", err)))
		return
	}
	log.Printf("DetailedSearch: Phase 2 finished - Found and preloaded data for %d companies.", len(paginatedData))

	// --- Этап 3: Формирование финального ответа (УПРОЩЕНО) ---
	// Нам больше не нужно создавать []models.CompanySearchResult
	// Переменная paginatedData уже содержит нужные данные []models.Registers с предзагруженными связями

	// Обновляем поле Data в ответе (хотя оно уже содержит данные после Find)
	response.Data = paginatedData

	// Лог перед финальным ответом
	log.Printf("DetailedSearch: Final check before response - TotalRecords: %d, Page: %d, Limit: %d, len(response.Data): %d",
		response.TotalRecords, response.Page, response.Limit, len(paginatedData))

	log.Printf("DetailedSearch: Request successful. Returning %d records.", len(paginatedData))
	c.JSON(http.StatusOK, response)
}

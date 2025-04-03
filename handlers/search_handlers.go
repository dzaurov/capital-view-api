// handlers/search_handlers.go
package handlers

import (
	"capital-view-api/db"     // Укажите правильный путь
	"capital-view-api/models" // Укажите правильный путь
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Вспомогательная структура только для получения ID из поиска
type RegCodeResult struct {
	Regcode                       *string
	LegalEntityRegistrationNumber *string // Добавим поле для owner/member
}

// DetailedSearch godoc
// @Summary Расширенный поиск компаний
// @Description Ищет компании по части названия, Regcode, SEPA, имени бенефициара или участника (member). Возвращает список компаний с полной информацией.
// @Tags search
// @Produce json
// @Param q query string true "Поисковый запрос (название, Regcode, SEPA, имя)"
// @Success 200 {array} models.CompanySearchResult "Список найденных компаний с детальной информацией"
// @Failure 400 {object} HTTPError "Неверный запрос - отсутствует параметр 'q'"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /search/detailed [get]
func DetailedSearch(c *gin.Context) {
	searchTerm := c.Query("q")
	if strings.TrimSpace(searchTerm) == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("поисковый параметр 'q' обязателен")))
		return
	}
	searchTermLike := "%" + searchTerm + "%"

	foundRegCodesMap := make(map[string]bool) // Используем map для уникальности Regcode

	var wg sync.WaitGroup    // Для параллельного поиска ID
	var mu sync.Mutex        // Для безопасной записи в map из горутин
	var searchErrors []error // Для сбора ошибок из горутин

	// --- Этап 1: Поиск совпадающих Regcode/LegalEntityRegistrationNumber ---

	// Поиск в 'register'
	wg.Add(1)
	go func() {
		defer wg.Done()
		var registerMatches []RegCodeResult
		err := db.DB.Model(&models.Register{}).
			Select("regcode"). // Выбираем только нужное поле
			Where("regcode = ?", searchTerm).
			Or("sepa = ?", searchTerm).
			Or("name LIKE ?", searchTermLike).
			Or("name_in_quotes LIKE ?", searchTermLike).
			Or("without_quotes LIKE ?", searchTermLike).
			Find(&registerMatches).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			mu.Lock()
			searchErrors = append(searchErrors, fmt.Errorf("ошибка поиска в register: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		for _, match := range registerMatches {
			if match.Regcode != nil && *match.Regcode != "" {
				foundRegCodesMap[*match.Regcode] = true
			}
		}
		mu.Unlock()
	}()

	// Поиск в 'beneficial_owner'
	wg.Add(1)
	go func() {
		defer wg.Done()
		var ownerMatches []RegCodeResult
		err := db.DB.Model(&models.BeneficialOwner{}).
			Select("legal_entity_registration_number"). // Выбираем только нужное поле
			Where("forename LIKE ?", searchTermLike).
			Or("surname LIKE ?", searchTermLike).
			Find(&ownerMatches).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			mu.Lock()
			searchErrors = append(searchErrors, fmt.Errorf("ошибка поиска в beneficial_owner: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		for _, match := range ownerMatches {
			// Используем поле LegalEntityRegistrationNumber из структуры RegCodeResult
			if match.LegalEntityRegistrationNumber != nil && *match.LegalEntityRegistrationNumber != "" {
				foundRegCodesMap[*match.LegalEntityRegistrationNumber] = true
			}
		}
		mu.Unlock()
	}()

	// Поиск в 'member'
	wg.Add(1)
	go func() {
		defer wg.Done()
		var memberMatches []RegCodeResult
		// Предполагаем, что ищем по основному номеру компании, к которой относится участник
		err := db.DB.Model(&models.Member{}).
			Select("legal_entity_registration_number"). // Выбираем только нужное поле
			Where("name LIKE ?", searchTermLike).
			// Возможно, нужно искать и по at_legal_entity_registration_number, если это ID компании, где он участник?
			// Or("at_legal_entity_registration_number = ?", searchTerm). // Если нужно искать и по этому полю
			Find(&memberMatches).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			mu.Lock()
			searchErrors = append(searchErrors, fmt.Errorf("ошибка поиска в member: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		for _, match := range memberMatches {
			// Используем поле LegalEntityRegistrationNumber из структуры RegCodeResult
			if match.LegalEntityRegistrationNumber != nil && *match.LegalEntityRegistrationNumber != "" {
				foundRegCodesMap[*match.LegalEntityRegistrationNumber] = true
			}
		}
		mu.Unlock()
	}()

	wg.Wait() // Дожидаемся завершения всех поисков ID

	// Проверяем ошибки поиска ID
	if len(searchErrors) > 0 {
		// Можно вернуть первую ошибку или скомбинировать их
		c.JSON(http.StatusInternalServerError, NewHTTPError(fmt.Errorf("ошибка при поиске ID компаний: %v", searchErrors)))
		return
	}

	if len(foundRegCodesMap) == 0 {
		c.JSON(http.StatusOK, []models.CompanySearchResult{}) // Возвращаем пустой массив, если ничего не найдено
		return
	}

	// --- Этап 2: Получение полной информации для каждого найденного Regcode ---

	finalResults := make([]models.CompanySearchResult, 0, len(foundRegCodesMap))
	var dataFetchErrors []error // Ошибки при получении полных данных

	// Получаем данные последовательно для каждого regCode
	// Можно оптимизировать с помощью горутин и каналов, если regCode много
	for regCode := range foundRegCodesMap {
		var companyRegister models.Register
		var members []models.Member
		var beneficialOwners []models.BeneficialOwner
		var financialStatements []models.FinancialStatement
		var reportDetails []models.FinancialReportDetail

		// Получаем основную информацию о компании
		err := db.DB.Where("regcode = ?", regCode).First(&companyRegister).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Компания была найдена по бенефициару/участнику, но нет записи в register? Странно, но пропустим.
				dataFetchErrors = append(dataFetchErrors, fmt.Errorf("не найдена запись register для regcode %s (ранее найденного)", regCode))
				continue // Переходим к следующему regCode
			} else {
				dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения register для %s: %w", regCode, err))
				continue
			}
		}

		// Получаем участников (Members)
		err = db.DB.Where("legal_entity_registration_number = ?", regCode).Find(&members).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения members для %s: %w", regCode, err))
			// Не прерываем, т.к. остальная информация может быть важна
		}

		// Получаем бенефициаров (Beneficial Owners)
		err = db.DB.Where("legal_entity_registration_number = ?", regCode).Find(&beneficialOwners).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения beneficial_owners для %s: %w", regCode, err))
			// Не прерываем
		}

		// Получаем финансовые отчеты (Financial Statements)
		err = db.DB.Where("legal_entity_registration_number = ?", regCode).Order("year desc").Find(&financialStatements).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения financial_statements для %s: %w", regCode, err))
			// Не прерываем
		} else if len(financialStatements) > 0 {
			// Получаем детали для каждого фин. отчета
			reportDetails = make([]models.FinancialReportDetail, 0, len(financialStatements))
			for _, fs := range financialStatements {
				// Копируем fs для избежания проблем с замыканием указателя в горутинах (если бы они были)
				currentFs := fs
				detail := models.FinancialReportDetail{
					FinancialStatementInfo: &currentFs,
				}
				financialStatementID := fs.ID // Используем ID из financial_statements

				// Используем GORM для поиска связанных записей по statement_id
				var incomeStatement models.IncomeStatement
				errIncome := db.DB.Where("statement_id = ?", financialStatementID).First(&incomeStatement).Error
				if errIncome == nil {
					detail.IncomeStatement = &incomeStatement
				} else if !errors.Is(errIncome, gorm.ErrRecordNotFound) {
					dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения income_statement для fs.id %d: %w", financialStatementID, errIncome))
				}

				var balanceSheet models.BalanceSheet
				errBalance := db.DB.Where("statement_id = ?", financialStatementID).First(&balanceSheet).Error
				if errBalance == nil {
					detail.BalanceSheet = &balanceSheet
				} else if !errors.Is(errBalance, gorm.ErrRecordNotFound) {
					dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения balance_sheet для fs.id %d: %w", financialStatementID, errBalance))
				}

				var cashFlowStatement models.CashFlowStatement
				errCashFlow := db.DB.Where("statement_id = ?", financialStatementID).First(&cashFlowStatement).Error
				if errCashFlow == nil {
					detail.CashFlowStatement = &cashFlowStatement
				} else if !errors.Is(errCashFlow, gorm.ErrRecordNotFound) {
					dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения cash_flow_statement для fs.id %d: %w", financialStatementID, errCashFlow))
				}

				reportDetails = append(reportDetails, detail)
			}
		}

		// Собираем результат для текущей компании
		companyResult := models.CompanySearchResult{
			RegisterInfo:     &companyRegister,
			Members:          members,
			BeneficialOwners: beneficialOwners,
			FinancialReports: reportDetails,
		}
		finalResults = append(finalResults, companyResult)
	}

	// Если были ошибки при получении данных, но есть и успешные результаты,
	// можно решить: вернуть частичные данные и залогировать ошибки,
	// или вернуть общую ошибку сервера. Вернем данные + залогируем ошибки.
	if len(dataFetchErrors) > 0 {
		// Логируем ошибки (используйте ваш логгер)
		log.Printf("Ошибки при получении полных данных для некоторых компаний: %v", dataFetchErrors)
	}

	// --- Этап 3: Возвращаем результат ---
	c.JSON(http.StatusOK, finalResults)
}

// Не забудьте определить HTTPError и NewHTTPError, если они еще не определены
// type HTTPError struct {
//    Error string `json:"error"`
// }
// func NewHTTPError(err error) HTTPError {
//    return HTTPError{Error: err.Error()}
// }

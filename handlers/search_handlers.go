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
	"time" // Добавим для логов

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Вспомогательная структура только для получения ID из поиска
type RegCodeResult struct {
	Regcode                       *string `gorm:"column:regcode"` // Добавим явное указание колонки (на всякий случай)
	LegalEntityRegistrationNumber *string `gorm:"column:legal_entity_registration_number"`
}

// DetailedSearch ... (остальные godoc комментарии) ...
func DetailedSearch(c *gin.Context) {
	startTime := time.Now()
	log.Println("DetailedSearch: Начало обработки запроса")

	searchTerm := c.Query("q")
	log.Printf("DetailedSearch: Получен поисковый запрос q='%s'", searchTerm)

	if strings.TrimSpace(searchTerm) == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("поисковый параметр 'q' обязателен")))
		return
	}

	// Преобразуем поисковый запрос в нижний регистр для поиска без учета регистра
	searchTermLower := strings.ToLower(searchTerm)
	searchTermLikeLower := "%" + searchTermLower + "%"
	log.Printf("DetailedSearch: Поиск будет регистронезависимым. Шаблон LIKE: '%s'", searchTermLikeLower)

	// !!! --- ВКЛЮЧАЕМ DEBUG GORM --- !!!
	// Создаем сессию с Debug() для логирования SQL
	debugDB := db.DB.Session(&gorm.Session{Logger: db.DB.Logger.LogMode(logger.Info)}) // <--- ИЗМЕНИТЬ log.Linfo на logger.Info // Используем стандартный логгер GORM
	// Если вы используете кастомный логгер для GORM, настройте его уровень на Info или Debug
	log.Println("DetailedSearch: GORM Debug режим включен для этого запроса.")

	foundRegCodesMap := make(map[string]bool) // Используем map для уникальности Regcode

	var wg sync.WaitGroup    // Для параллельного поиска ID
	var mu sync.Mutex        // Для безопасной записи в map из горутин
	var searchErrors []error // Для сбора ошибок из горутин

	log.Println("DetailedSearch: Этап 1: Запуск параллельного поиска совпадающих ID")

	// --- Этап 1: Поиск совпадающих Regcode/LegalEntityRegistrationNumber ---

	// Поиск в 'register' (с LOWER и Debug)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var registerMatches []RegCodeResult
		log.Println("DetailedSearch: Goroutine(register): Начало поиска...")
		// Используем debugDB и LOWER() для регистронезависимого поиска
		err := debugDB.Model(&models.Register{}). // Используем debugDB
								Select("regcode").
								Where("LOWER(regcode) = ?", searchTermLower). // Поиск по regcode тоже сделаем регистронезависимым
								Or("LOWER(sepa) = ?", searchTermLower).       // И по sepa
								Or("LOWER(name) LIKE ?", searchTermLikeLower).
								Or("LOWER(name_in_quotes) LIKE ?", searchTermLikeLower).
								Or("LOWER(without_quotes) LIKE ?", searchTermLikeLower).
								Find(&registerMatches).Error

		// --- ЛОГИРОВАНИЕ СРАЗУ ПОСЛЕ ЗАПРОСА ---
		if err != nil {
			log.Printf("DetailedSearch: Goroutine(register): Ошибка выполнения запроса: %v", err)
			// Не выходим сразу при gorm.ErrRecordNotFound, просто логируем
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				mu.Lock()
				searchErrors = append(searchErrors, fmt.Errorf("ошибка поиска в register: %w", err))
				mu.Unlock()
			}
			// Логируем количество найденных даже при ошибке (если Find что-то вернул до ошибки)
			log.Printf("DetailedSearch: Goroutine(register): Найдено совпадений (до обработки): %d", len(registerMatches))
		} else {
			log.Printf("DetailedSearch: Goroutine(register): Запрос успешно выполнен. Найдено совпадений (до обработки): %d", len(registerMatches))
		}
		// --- КОНЕЦ ЛОГИРОВАНИЯ ---

		addedCount := 0
		mu.Lock()
		for i, match := range registerMatches {
			// Дополнительное логирование самих найденных значений
			log.Printf("DetailedSearch: Goroutine(register): Обработка совпадения %d: Regcode=%v", i+1, match.Regcode)
			if match.Regcode != nil && *match.Regcode != "" {
				if !foundRegCodesMap[*match.Regcode] {
					foundRegCodesMap[*match.Regcode] = true
					addedCount++
				}
			}
		}
		mu.Unlock()
		log.Printf("DetailedSearch: Goroutine(register): Завершение обработки. Добавлено новых уникальных Regcode: %d", addedCount)
	}()

	// Поиск в 'beneficial_owner' (с LOWER и Debug)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var ownerMatches []RegCodeResult
		log.Println("DetailedSearch: Goroutine(beneficial_owner): Начало поиска...")
		// Используем debugDB и LOWER()
		err := debugDB.Model(&models.BeneficialOwner{}). // Используем debugDB
									Select("legal_entity_registration_number").
									Where("LOWER(forename) LIKE ?", searchTermLikeLower). // Ищем регистронезависимо
									Or("LOWER(surname) LIKE ?", searchTermLikeLower).
									Find(&ownerMatches).Error

		// --- ЛОГИРОВАНИЕ СРАЗУ ПОСЛЕ ЗАПРОСА ---
		if err != nil {
			log.Printf("DetailedSearch: Goroutine(beneficial_owner): Ошибка выполнения запроса: %v", err)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				mu.Lock()
				searchErrors = append(searchErrors, fmt.Errorf("ошибка поиска в beneficial_owner: %w", err))
				mu.Unlock()
			}
			log.Printf("DetailedSearch: Goroutine(beneficial_owner): Найдено совпадений (до обработки): %d", len(ownerMatches))
		} else {
			log.Printf("DetailedSearch: Goroutine(beneficial_owner): Запрос успешно выполнен. Найдено совпадений (до обработки): %d", len(ownerMatches))
		}
		// --- КОНЕЦ ЛОГИРОВАНИЯ ---

		addedCount := 0
		mu.Lock()
		for i, match := range ownerMatches {
			log.Printf("DetailedSearch: Goroutine(beneficial_owner): Обработка совпадения %d: LegalEntityRegistrationNumber=%v", i+1, match.LegalEntityRegistrationNumber)
			if match.LegalEntityRegistrationNumber != nil && *match.LegalEntityRegistrationNumber != "" {
				if !foundRegCodesMap[*match.LegalEntityRegistrationNumber] {
					foundRegCodesMap[*match.LegalEntityRegistrationNumber] = true
					addedCount++
				}
			}
		}
		mu.Unlock()
		log.Printf("DetailedSearch: Goroutine(beneficial_owner): Завершение обработки. Добавлено новых уникальных ID: %d", addedCount)
	}()

	// Поиск в 'member' (с LOWER и Debug)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var memberMatches []RegCodeResult
		log.Println("DetailedSearch: Goroutine(member): Начало поиска...")
		// Используем debugDB и LOWER()
		err := debugDB.Model(&models.Member{}). // Используем debugDB
							Select("legal_entity_registration_number").
							Where("LOWER(name) LIKE ?", searchTermLikeLower). // Ищем регистронезависимо
			// Or("at_legal_entity_registration_number = ?", searchTerm). // Если нужно - тоже добавить LOWER()
			Find(&memberMatches).Error

		// --- ЛОГИРОВАНИЕ СРАЗУ ПОСЛЕ ЗАПРОСА ---
		if err != nil {
			log.Printf("DetailedSearch: Goroutine(member): Ошибка выполнения запроса: %v", err)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				mu.Lock()
				searchErrors = append(searchErrors, fmt.Errorf("ошибка поиска в member: %w", err))
				mu.Unlock()
			}
			log.Printf("DetailedSearch: Goroutine(member): Найдено совпадений (до обработки): %d", len(memberMatches))
		} else {
			log.Printf("DetailedSearch: Goroutine(member): Запрос успешно выполнен. Найдено совпадений (до обработки): %d", len(memberMatches))
		}
		// --- КОНЕЦ ЛОГИРОВАНИЯ ---

		addedCount := 0
		mu.Lock()
		for i, match := range memberMatches {
			log.Printf("DetailedSearch: Goroutine(member): Обработка совпадения %d: LegalEntityRegistrationNumber=%v", i+1, match.LegalEntityRegistrationNumber)
			if match.LegalEntityRegistrationNumber != nil && *match.LegalEntityRegistrationNumber != "" {
				if !foundRegCodesMap[*match.LegalEntityRegistrationNumber] {
					foundRegCodesMap[*match.LegalEntityRegistrationNumber] = true
					addedCount++
				}
			}
		}
		mu.Unlock()
		log.Printf("DetailedSearch: Goroutine(member): Завершение обработки. Добавлено новых уникальных ID: %d", addedCount)
	}()

	log.Println("DetailedSearch: Этап 1: Ожидание завершения всех горутин поиска ID...")
	wg.Wait() // Дожидаемся завершения всех поисков ID
	log.Println("DetailedSearch: Этап 1: Все горутины поиска ID завершены.")

	// Проверяем ошибки поиска ID
	if len(searchErrors) > 0 {
		log.Printf("DetailedSearch: Ошибка: Обнаружены ошибки во время Этапа 1 (поиск ID): %v", searchErrors)
		c.JSON(http.StatusInternalServerError, NewHTTPError(fmt.Errorf("ошибка при поиске ID компаний: %v", searchErrors)))
		return
	}

	log.Printf("DetailedSearch: Этап 1: Завершен успешно. Найдено уникальных ID в foundRegCodesMap: %d", len(foundRegCodesMap))

	// --- ДАЛЬШЕ КОД ОСТАЕТСЯ БЕЗ ИЗМЕНЕНИЙ ---

	if len(foundRegCodesMap) == 0 {
		log.Println("DetailedSearch: Уникальные ID не найдены в foundRegCodesMap. Возвращаем пустой результат.") // Добавим лог сюда
		c.JSON(http.StatusOK, []models.CompanySearchResult{})                                                    // Возвращаем пустой массив, если ничего не найдено
		return
	}

	// --- Этап 2: Получение полной информации ... (код без изменений) ---
	log.Println("DetailedSearch: Этап 2: Начало получения полной информации для найденных ID.")
	finalResults := make([]models.CompanySearchResult, 0, len(foundRegCodesMap))
	var dataFetchErrors []error

	// Используем обычный db.DB здесь, Debug не обязателен для всех запросов
	for regCode := range foundRegCodesMap {
		log.Printf("DetailedSearch: Этап 2: Обработка ID: %s", regCode)
		// ... (весь ваш код для Этапа 2 без изменений) ...

		var companyRegister models.Register
		var members []models.Member
		var beneficialOwners []models.BeneficialOwner
		var financialStatements []models.FinancialStatement
		var reportDetails []models.FinancialReportDetail

		// Получаем основную информацию о компании
		err := db.DB.Where("regcode = ?", regCode).First(&companyRegister).Error // Используем db.DB
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("DetailedSearch: Этап 2: Предупреждение: Не найдена запись register для regcode %s. Пропускаем.", regCode)
				dataFetchErrors = append(dataFetchErrors, fmt.Errorf("не найдена запись register для regcode %s (ранее найденного)", regCode))
				continue // Переходим к следующему regCode
			} else {
				log.Printf("DetailedSearch: Этап 2: Ошибка получения register для %s: %v. Пропускаем.", regCode, err)
				dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения register для %s: %w", regCode, err))
				continue
			}
		}

		// ... (остальной код получения members, beneficialOwners, financialStatements...) ...
		// Получаем участников (Members)
		err = db.DB.Where("legal_entity_registration_number = ?", regCode).Find(&members).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("DetailedSearch: Этап 2: Ошибка получения members для %s: %v", regCode, err)
			dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения members для %s: %w", regCode, err))
		}

		// Получаем бенефициаров (Beneficial Owners)
		err = db.DB.Where("legal_entity_registration_number = ?", regCode).Find(&beneficialOwners).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("DetailedSearch: Этап 2: Ошибка получения beneficial_owners для %s: %v", regCode, err)
			dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения beneficial_owners для %s: %w", regCode, err))
		}

		// Получаем финансовые отчеты (Financial Statements)
		err = db.DB.Where("legal_entity_registration_number = ?", regCode).Order("year desc").Find(&financialStatements).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("DetailedSearch: Этап 2: Ошибка получения financial_statements для %s: %v", regCode, err)
			dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения financial_statements для %s: %w", regCode, err))
		} else if len(financialStatements) > 0 {
			// Получаем детали для каждого фин. отчета
			reportDetails = make([]models.FinancialReportDetail, 0, len(financialStatements))
			for _, fs := range financialStatements {
				currentFs := fs
				detail := models.FinancialReportDetail{
					FinancialStatementInfo: &currentFs,
				}
				financialStatementID := fs.ID

				var incomeStatement models.IncomeStatement
				errIncome := db.DB.Where("statement_id = ?", financialStatementID).First(&incomeStatement).Error
				if errIncome == nil {
					detail.IncomeStatement = &incomeStatement
				} else if !errors.Is(errIncome, gorm.ErrRecordNotFound) {
					log.Printf("DetailedSearch: Этап 2: Ошибка получения income_statement для fs.id %d: %v", financialStatementID, errIncome)
					dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения income_statement для fs.id %d: %w", financialStatementID, errIncome))
				}

				var balanceSheet models.BalanceSheet
				errBalance := db.DB.Where("statement_id = ?", financialStatementID).First(&balanceSheet).Error
				if errBalance == nil {
					detail.BalanceSheet = &balanceSheet
				} else if !errors.Is(errBalance, gorm.ErrRecordNotFound) {
					log.Printf("DetailedSearch: Этап 2: Ошибка получения balance_sheet для fs.id %d: %v", financialStatementID, errBalance)
					dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения balance_sheet для fs.id %d: %w", financialStatementID, errBalance))
				}

				var cashFlowStatement models.CashFlowStatement
				errCashFlow := db.DB.Where("statement_id = ?", financialStatementID).First(&cashFlowStatement).Error
				if errCashFlow == nil {
					detail.CashFlowStatement = &cashFlowStatement
				} else if !errors.Is(errCashFlow, gorm.ErrRecordNotFound) {
					log.Printf("DetailedSearch: Этап 2: Ошибка получения cash_flow_statement для fs.id %d: %v", financialStatementID, errCashFlow)
					dataFetchErrors = append(dataFetchErrors, fmt.Errorf("ошибка получения cash_flow_statement для fs.id %d: %w", financialStatementID, errCashFlow))
				}

				reportDetails = append(reportDetails, detail)
			}
		}

		companyResult := models.CompanySearchResult{
			RegisterInfo:     &companyRegister,
			Members:          members,
			BeneficialOwners: beneficialOwners,
			FinancialReports: reportDetails,
		}
		finalResults = append(finalResults, companyResult)
		log.Printf("DetailedSearch: Этап 2: Результат для ID %s собран.", regCode)
	} // Конец цикла for regCode

	log.Printf("DetailedSearch: Этап 2: Завершен. Собрано %d полных результатов.", len(finalResults))

	if len(dataFetchErrors) > 0 {
		log.Printf("DetailedSearch: Обнаружены ошибки (%d) во время Этапа 2 (сбор полных данных): %v", len(dataFetchErrors), dataFetchErrors)
	}

	// --- Этап 3: Возвращаем результат ---
	duration := time.Since(startTime)
	log.Printf("DetailedSearch: Отправка итогового результата клиенту (%d записей). Общее время: %s", len(finalResults), duration)
	c.JSON(http.StatusOK, finalResults)
}

// --- Не забудьте про эту функцию ---
// Убедитесь, что она определена ОДИН РАЗ в вашем пакете handlers
// func NewHTTPError(err error) map[string]interface{} {
// 	 return gin.H{"error": err.Error()}
// }

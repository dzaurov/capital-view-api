// cmd/importer/main.go
package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	// Используем алиас для пакета db вашего проекта
	dbConn "capital-view-api/db" // <--- Проверьте правильность пути
	"capital-view-api/models"    // <--- Проверьте правильность пути

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Config для обработки CSV
type Config struct {
	FileName       string          // Имя CSV файла (без расширения)
	Model          interface{}     // Указатель на пустую структуру модели (например, &models.Registers{})
	ConflictTarget []clause.Column // Колонки для ON CONFLICT
	UpdateColumns  []string        // Колонки для UPDATE в ON CONFLICT
}

func main() {
	// --- Настройка ---
	csvDir := flag.String("csvdir", "./csv_real", "Directory containing CSV files")
	flag.Parse()
	log.Printf("Starting CSV import from directory: %s", *csvDir)

	// --- Подключение к БД ---
	log.Println("Connecting to database...")
	if err := dbConn.ConnectDatabase(); err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	db := dbConn.DB
	log.Println("Database connection successful.")

	// --- AutoMigrate ---
	log.Println("Running AutoMigrate...")
	err := db.AutoMigrate( // Добавляем ВСЕ модели для создания/обновления таблиц и ИНДЕКСОВ
		&models.Registers{}, // <-- Используем новую модель
		&models.Member{},
		&models.BeneficialOwner{},
		&models.FinancialStatement{},
		&models.IncomeStatement{},
		&models.BalanceSheet{},
		&models.CashFlowStatement{},
		// Добавьте сюда &models.Officer{} если будете импортировать officers.csv
	)
	if err != nil {
		log.Fatalf("FATAL: AutoMigrate failed: %v", err)
	}
	log.Println("AutoMigrate completed.")

	// --- Конфигурация импорта (СКОРРЕКТИРОВАНА!) ---
	configs := map[string]Config{
		"registers": {
			FileName:       "register",          // CSV файл называется register.csv
			Model:          &models.Registers{}, // Используем модель Registers (для таблицы registers)
			ConflictTarget: []clause.Column{{Name: "regcode"}},
			UpdateColumns: []string{ // Колонки из models.Registers
				"sepa", "name", "name_before_quotes", "name_in_quotes", "name_after_quotes",
				"without_quotes", "regtype", "regtype_text", "type", "type_text", "registered",
				"terminated", "closed", "address", "index_company", // <-- index_company
				"addressid", "region", "city", "atvk", "reregistration_term",
			},
		},
		"members": {
			FileName:       "members", // CSV файл members.csv
			Model:          &models.Member{},
			ConflictTarget: []clause.Column{{Name: "id"}}, // Используем ID из CSV
			UpdateColumns: []string{ // Все поля модели Member КРОМЕ ID
				"uri", "at_legal_entity_registration_number", "entity_type", "name",
				"latvian_identity_number_masked", "birth_date", "legal_entity_registration_number",
				"number_of_shares", "share_nominal_value", "share_currency", "date_from",
				"registered_on", "last_modified_at",
			},
		},
		"beneficial_owners": {
			FileName:       "beneficial_owners", // CSV файл beneficial_owners.csv
			Model:          &models.BeneficialOwner{},
			ConflictTarget: []clause.Column{{Name: "id"}}, // Используем ID из CSV
			UpdateColumns: []string{ // Все поля модели BeneficialOwner КРОМЕ ID
				"legal_entity_registration_number", "forename", "surname",
				"latvian_identity_number_masked", "birth_date", "nationality", "residence",
				"registered_on", "last_modified_at",
			},
		},
		"financial_statements": {
			FileName: "financial_statements", // CSV файл financial_statements.csv
			Model:    &models.FinancialStatement{},
			ConflictTarget: []clause.Column{ // Ключ: компания + год
				{Name: "legal_entity_registration_number"},
				{Name: "year"},
			},
			UpdateColumns: []string{ // Все поля КРОМЕ ID и ключей конфликта
				"file_id", "source_schema", "source_type", "year_started_on", "year_ended_on",
				"employees", "rounded_to_nearest", "currency", "created_at",
			},
		},
		"income_statements": {
			FileName:       "income_statements", // CSV файл income_statements.csv
			Model:          &models.IncomeStatement{},
			ConflictTarget: []clause.Column{{Name: "statement_id"}}, // Ключ (теперь с uniqueIndex)
			UpdateColumns: []string{ // Все поля КРОМЕ ID и statement_id
				"file_id", "net_turnover", "by_nature_inventory_change", "by_nature_long_term_investment_expenses",
				"by_nature_other_operating_revenues", "by_nature_material_expenses", "by_nature_labour_expenses",
				"by_nature_depreciation_expenses", "by_function_cost_of_goods_sold", "by_function_gross_profit",
				"by_function_selling_expenses", "by_function_administrative_expenses",
				"by_function_other_operating_revenues", "other_operating_expenses", "equity_investment_earnings",
				"other_long_term_investment_earnings", "other_interest_revenues", "investment_fair_value_adjustments",
				"interest_expenses", "extra_revenues", "extra_expenses", "income_before_income_taxes",
				"provision_for_income_taxes", "income_after_income_taxes", "other_taxes", "extra_dividends", "net_income",
			},
		},
		"balance_sheets": {
			FileName:       "balance_sheets", // CSV файл balance_sheets.csv
			Model:          &models.BalanceSheet{},
			ConflictTarget: []clause.Column{{Name: "statement_id"}}, // Ключ (теперь с uniqueIndex)
			UpdateColumns: []string{ // Все поля КРОМЕ ID и statement_id
				"file_id", "cash", "marketable_securities", "accounts_receivable", "inventories",
				"total_current_assets", "investments", "fixed_assets", "intangible_assets",
				"total_non_current_assets", "total_assets", "future_housing_repairs_payments",
				"current_liabilities", "non_current_liabilities", "provisions", "equity", "total_equities",
			},
		},
		"cash_flow_statements": {
			FileName:       "cash_flow_statements", // CSV файл cash_flow_statements.csv
			Model:          &models.CashFlowStatement{},
			ConflictTarget: []clause.Column{{Name: "statement_id"}}, // Ключ (теперь с uniqueIndex)
			UpdateColumns: []string{ // Все поля КРОМЕ ID и statement_id
				"file_id", "cfo_dm_cash_received_from_customers", "cfo_dm_cash_paid_to_suppliers_employees",
				"cfo_dm_other_cash_received_paid", "cfo_dm_operating_cash_flow", "cfo_dm_interest_paid",
				"cfo_dm_income_taxes_paid", "cfo_dm_extra_items_cash_flow", "cfo_dm_net_operating_cash_flow",
				"cfo_im_income_before_income_taxes", "cfo_im_income_before_changes_in_working_capital",
				"cfo_im_operating_cash_flow", "cfo_im_interest_paid", "cfo_im_income_taxes_paid",
				"cfo_im_extra_items_cash_flow", "cfo_im_net_operating_cash_flow", "cfi_acquisition_of_stocks_shares",
				"cfi_sale_proceeds_from_stocks_shares", "cfi_acquisition_of_fixed_assets_intangible_assets",
				"cfi_sale_proceeds_from_fixed_assets_intangible_assets", "cfi_loans_made",
				"cfi_repayments_of_loans_received", "cfi_interest_received", "cfi_dividends_received",
				"cfi_net_investing_cash_flow", "cff_proceeds_from_stocks_bonds_issuance_or_contributed_capital",
				"cff_loans_received", "cff_subsidies_grants_donations_received", "cff_repayments_of_loans_made",
				"cff_repayments_of_lease_obligations", "cff_dividends_paid", "cff_net_financing_cash_flow",
				"effect_of_exchange_rate_change", "net_increase", "at_beginning_of_year", "at_end_of_year",
			},
		},
		// Добавьте сюда officers, если нужно, создав модель и конфигурацию
	}

	// --- Обработка конфигураций ---
	ctx := context.Background()
	for cfgName, cfg := range configs {
		// Проверяем путь к директории CSV
		if _, err := os.Stat(*csvDir); os.IsNotExist(err) {
			log.Printf("WARN: CSV directory not found: %s, skipping all files.", *csvDir)
			break // Выходим, если папки нет
		}
		filePath := filepath.Join(*csvDir, cfg.FileName+".csv")
		log.Printf("Processing file: %s for table %s", filePath, cfgName)
		err := processCSV(ctx, db, filePath, cfg)
		if err != nil {
			log.Printf("ERROR processing %s: %v", filePath, err)
		} else {
			log.Printf("Successfully finished processing file: %s", filePath)
		}
	}
	log.Println("CSV import process finished.")
}

// processCSV обрабатывает один CSV файл
func processCSV(ctx context.Context, db *gorm.DB, filePath string, cfg Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("WARN: File %s not found, skipping.", filePath)
			return nil // Не ошибка, если файл не найден
		}
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // <--- УСТАНОВИТЕ ПРАВИЛЬНЫЙ РАЗДЕЛИТЕЛЬ (',' или ';')
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err == io.EOF {
		log.Printf("WARN: File %s is empty.", filePath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not read header: %w", err)
	}
	log.Printf("Headers found: %d", len(headers))

	headerMap := make(map[string]int, len(headers))
	for i, h := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Получаем схему модели один раз перед циклом
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(cfg.Model); err != nil {
		return fmt.Errorf("could not parse model schema %T: %w", cfg.Model, err)
	}
	schema := stmt.Schema

	recordsProcessed := 0
	recordsUpserted := 0
	recordsFailed := 0
	modelType := reflect.TypeOf(cfg.Model).Elem() // Тип структуры (не указателя)

	// Используем транзакцию для каждого файла для ускорения
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction for file %s: %w", filePath, tx.Error)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("WARN: Error reading row (approx line %d): %v - Skipping row", recordsProcessed+2, err)
			recordsFailed++
			continue
		}

		recordsProcessed++

		currentRecord := reflect.New(modelType).Interface() // Создаем *указатель* на структуру
		currentRecordValue := reflect.ValueOf(currentRecord)

		parseErrorsInRow := false
		for _, field := range schema.Fields {
			// Пропускаем ID, если он автоинкрементный (мы используем его только как ConflictTarget из CSV)
			if field.Name == "ID" && field.AutoIncrement && !isConflictTarget(field.Name, cfg.ConflictTarget) {
				continue
			}

			// --- Маппинг для index -> index_company ---
			csvHeaderName := strings.ToLower(field.DBName)
			// Специальный случай для register.csv -> registers
			if schema.Table == "registers" && csvHeaderName == "index_company" {
				csvHeaderName = "index" // Ищем "index" в заголовке CSV
			}
			// -----------------------------------------

			columnIndex, headerFound := headerMap[csvHeaderName]
			if !headerFound {
				// Поле модели не найдено в CSV - пропускаем
				// log.Printf("DEBUG: Field %s (DB: %s / CSV: %s) not found in CSV header.", field.Name, field.DBName, csvHeaderName)
				continue
			}
			if columnIndex >= len(row) {
				log.Printf("WARN: Row %d is shorter than expected header index for '%s'. Skipping field.", recordsProcessed+1, csvHeaderName)
				parseErrorsInRow = true // Отмечаем ошибку
				continue
			}

			valueStr := strings.TrimSpace(row[columnIndex])

			err := SetFieldValue(ctx, currentRecordValue, field, valueStr)
			if err != nil {
				log.Printf("WARN: Failed to set field '%s' from value '%s' (line %d): %v", field.Name, valueStr, recordsProcessed+1, err)
				parseErrorsInRow = true
			}
		}

		if parseErrorsInRow {
			recordsFailed++
			log.Printf("WARN: Skipping row %d due to parsing errors.", recordsProcessed+1)
			continue
		}

		// Выполняем Upsert внутри транзакции
		result := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   cfg.ConflictTarget,
			DoUpdates: clause.AssignmentColumns(cfg.UpdateColumns),
		}).Create(currentRecord)

		if result.Error != nil {
			log.Printf("ERROR: Failed to upsert record (line %d): %v - Record Data approx: %+v", recordsProcessed+1, result.Error, currentRecord)
			recordsFailed++
			// При ошибке в транзакции, откатываем всю транзакцию для файла и выходим из функции
			tx.Rollback()
			return fmt.Errorf("failed to upsert record (line %d), transaction rolled back: %w", recordsProcessed+1, result.Error)
		} else {
			if result.RowsAffected > 0 {
				recordsUpserted++
			}
			if recordsProcessed%1000 == 0 { // Логируем прогресс
				log.Printf("Processed %d rows for %s...", recordsProcessed, filePath)
			}
		}
	}

	// Коммитим транзакцию, если весь файл обработан без ошибок upsert
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction for file %s: %v", filePath, err)
		// Ошибки чтения CSV уже обработаны, но ошибка Commit критична
		recordsFailed = recordsProcessed - recordsUpserted // Считаем все непрошедшие строки как ошибки
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Finished processing %s. Total rows read: %d, Rows upserted/updated: %d, Rows failed/skipped: %d",
		filePath, recordsProcessed, recordsUpserted, recordsFailed)

	return nil
}

// isConflictTarget проверяет, является ли поле частью ключа конфликта
func isConflictTarget(fieldName string, targets []clause.Column) bool {
	for _, target := range targets {
		if target.Name == fieldName {
			return true
		}
	}
	return false
}

// SetFieldValue - Хелпер для установки значения поля структуры через reflect
func SetFieldValue(ctx context.Context, targetStructValue reflect.Value, field *schema.Field, valueStr string) error {

	// ID обрабатываем особо
	if field.Name == "ID" {
		if valueStr == "" {
			if field.AutoIncrement {
				return nil
			} // Пропускаем автоинкрементный ID, если пуст
			return fmt.Errorf("ID field (PK) cannot be empty in CSV")
		}
		idVal, err := strconv.ParseUint(valueStr, 10, 0)
		if err != nil {
			return fmt.Errorf("parsing ID '%s' as uint failed: %w", valueStr, err)
		}
		// Устанавливаем ID правильного типа (uint)
		return field.Set(ctx, targetStructValue, reflect.ValueOf(uint(idVal)).Convert(field.FieldType).Interface())
	}

	// Проверяем на пустую строку для остальных полей
	if valueStr == "" {
		if field.FieldType.Kind() == reflect.Ptr {
			nilValue := reflect.Zero(field.FieldType) // typed nil pointer (*string)
			return field.Set(ctx, targetStructValue, nilValue.Interface())
		} else {
			// Если поле не указатель, но строка пустая - оставляем zero value
			return nil
		}
	}

	// --- Все поля, кроме ID, у нас *string согласно моделям ---
	// Проверяем, что поле действительно *string
	if field.FieldType.Kind() == reflect.Ptr && field.FieldType.Elem().Kind() == reflect.String {
		// Создаем указатель на строку
		ptrValue := reflect.New(field.FieldType.Elem())                // *string -> string
		ptrValue.Elem().SetString(valueStr)                            // Устанавливаем значение строки
		return field.Set(ctx, targetStructValue, ptrValue.Interface()) // Устанавливаем указатель *string
	} else {
		// Если поле в модели не *string (кроме ID), это несоответствие
		// или нужно добавить логику парсинга для других типов здесь
		log.Printf("DEBUG: Field %s is not *string (and not ID), value '%s' ignored.", field.Name, valueStr)
		// Возвращаем ошибку или просто игнорируем? Давайте игнорировать пока.
		// return fmt.Errorf("field %s type mismatch: expected *string, got %s", field.Name, field.FieldType.String())
		return nil
	}

	// Старая логика парсинга чисел/bool закомментирована, т.к. модели используют *string
	/*
		var value interface{}
		var err error
		fieldType := field.FieldType
		isPtr := false
		if fieldType.Kind() == reflect.Ptr {	fieldType = fieldType.Elem(); isPtr = true }

		switch fieldType.Kind() {
		// ... парсинг uint, int, float, bool ...
		}
		if err != nil { return fmt.Errorf(...) }
		if isPtr { // ... установка указателя ... } else { // ... установка значения ... }
	*/
}

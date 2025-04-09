package main

import (
	"capital-view-api/db"       // Adjust import path if needed
	_ "capital-view-api/docs"   // Adjust import path (important for swag init)
	"capital-view-api/handlers" // Adjust import path if needed
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Your API Title
// @version 1.0
// @description This is a sample server for managing data based on the provided schema.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	// Initialize Database
	err := db.ConnectDatabase() // <--- Используйте ConnectDatabase() вместо InitDB()
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	// --- !!! ЯВНОЕ СОЗДАНИЕ УНИКАЛЬНЫХ И ОБЫЧНЫХ ИНДЕКСОВ !!! ---
	log.Println("Ensuring necessary indexes exist...")
	indexCommands := []string{
		// Уникальные
		"CREATE UNIQUE INDEX IF NOT EXISTS uq_fs_company_year ON financial_statements (legal_entity_registration_number, year)",
		"CREATE UNIQUE INDEX IF NOT EXISTS uq_is_statement_id ON income_statements (statement_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS uq_bs_statement_id ON balance_sheets (statement_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS uq_cfs_statement_id ON cash_flow_statements (statement_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS uq_reg_regcode ON registers (regcode)",

		// Обычные индексы для ускорения поиска LIKE и связей Preload
		"CREATE INDEX IF NOT EXISTS idx_registers_name ON registers (name)",                                     // <-- Для LIKE
		"CREATE INDEX IF NOT EXISTS idx_registers_name_in_quotes ON registers (name_in_quotes)",                 // <-- Для LIKE
		"CREATE INDEX IF NOT EXISTS idx_registers_without_quotes ON registers (without_quotes)",                 // <-- Для LIKE
		"CREATE INDEX IF NOT EXISTS idx_members_regcode ON members (legal_entity_registration_number)",          // <-- Для Preload
		"CREATE INDEX IF NOT EXISTS idx_members_name ON members (name)",                                         // <-- Для LIKE
		"CREATE INDEX IF NOT EXISTS idx_owners_regcode ON beneficial_owners (legal_entity_registration_number)", // <-- Для Preload
		"CREATE INDEX IF NOT EXISTS idx_owners_forename ON beneficial_owners (forename)",                        // <-- Для LIKE
		"CREATE INDEX IF NOT EXISTS idx_owners_surname ON beneficial_owners (surname)",                          // <-- Для LIKE
		// Индекс для financial_statements.legal_entity_registration_number уже покрыт уникальным составным
	}
	for _, cmd := range indexCommands {
		if tx := db.DB.Exec(cmd); tx.Error != nil {
			log.Printf("WARN: Failed to execute index command: %v - SQL: %s", tx.Error, cmd)
		}
	}
	log.Println("Indexes checked/created.")
	// ---------------------------------------------------------

	// Initialize Gin router
	router := gin.Default()

	// Настройка эндпоинтов API v1
	v1 := router.Group("/api/v1")
	{
		// --- !!! ДОБАВИТЬ РОУТ ДЛЯ ПОЛНОЙ ИНФОРМАЦИИ О КОМПАНИИ !!! ---
		v1.GET("/company/:regcode", handlers.GetCompanyDetailsByRegcode)
		// ---------------------------------------------------------------

		// Register routes
		v1.GET("/registers", handlers.GetAllRegisters)
		v1.GET("/register/:regcode", handlers.GetRegisterByID) // Этот остается для базовой инфы
		// ... (закомментированные CRUD роуты) ...

		// Member routes
		v1.GET("/members/by-regcode/:regcode", handlers.GetMembersByRegcode)
		// ... (закомментированные CRUD роуты) ...

		// Beneficial Owner routes
		v1.GET("/beneficial-owners/by-regcode/:regcode", handlers.GetBeneficialOwnersByRegcode)
		// ... (закомментированные CRUD роуты) ...

		// Financial Statement routes
		v1.GET("/financial-statements/by-regcode/:regcode", handlers.GetFinancialStatementsByRegcode)
		// ... (закомментированные CRUD роуты) ...

		// Search routes
		searchGroup := v1.Group("/search")
		{
			// Этот роут теперь возвращает УПРОЩЕННЫЕ данные
			searchGroup.GET("/detailed", handlers.DetailedSearch)
		}
	} // конец v1

	// Swagger Documentation Route
	// The URL will be http://localhost:8080/swagger/index.html (or your host/port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	port := ":8080" // You can make this configurable
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

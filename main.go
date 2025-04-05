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

	// Initialize Gin router
	router := gin.Default()

	// API v1 Group
	v1 := router.Group("/api/v1")
	{
		// --- Register Routes ---
		registers := v1.Group("/registers")
		{
			registers.POST("", handlers.CreateRegister)      // POST /api/v1/registers
			registers.GET("", handlers.GetRegisters)         // GET /api/v1/registers
			registers.GET(":id", handlers.GetRegister)       // GET /api/v1/registers/:id
			registers.PUT(":id", handlers.UpdateRegister)    // PUT /api/v1/registers/:id
			registers.DELETE(":id", handlers.DeleteRegister) // DELETE /api/v1/registers/:id
		}

		// --- Member Routes ---
		members := v1.Group("/members")
		{
			// *** UNCOMMENT AND ADD THESE LINES ***
			members.POST("", handlers.CreateMember)      // POST /api/v1/members
			members.GET("", handlers.GetMembers)         // GET /api/v1/members
			members.GET(":id", handlers.GetMember)       // GET /api/v1/members/:id
			members.PUT(":id", handlers.UpdateMember)    // PUT /api/v1/members/:id
			members.DELETE(":id", handlers.DeleteMember) // DELETE /api/v1/members/:id
		}

		// --- Income Statement Routes ---
		incomeStatements := v1.Group("/income-statements")
		{
			// *** UNCOMMENT AND ADD THESE LINES ***
			incomeStatements.POST("", handlers.CreateIncomeStatement)      // POST /api/v1/income-statements
			incomeStatements.GET("", handlers.GetIncomeStatements)         // GET /api/v1/income-statements
			incomeStatements.GET(":id", handlers.GetIncomeStatement)       // GET /api/v1/income-statements/:id
			incomeStatements.PUT(":id", handlers.UpdateIncomeStatement)    // PUT /api/v1/income-statements/:id
			incomeStatements.DELETE(":id", handlers.DeleteIncomeStatement) // DELETE /api/v1/income-statements/:id
		}

		// --- Financial Statement Routes ---
		financialStatements := v1.Group("/financial-statements")
		{
			// *** UNCOMMENT AND ADD THESE LINES ***
			financialStatements.POST("", handlers.CreateFinancialStatement)      // POST /api/v1/financial-statements
			financialStatements.GET("", handlers.GetFinancialStatements)         // GET /api/v1/financial-statements
			financialStatements.GET(":id", handlers.GetFinancialStatement)       // GET /api/v1/financial-statements/:id
			financialStatements.PUT(":id", handlers.UpdateFinancialStatement)    // PUT /api/v1/financial-statements/:id
			financialStatements.DELETE(":id", handlers.DeleteFinancialStatement) // DELETE /api/v1/financial-statements/:id
		}

		// --- Cash Flow Statement Routes ---
		cashFlowStatements := v1.Group("/cash-flow-statements")
		{
			// *** UNCOMMENT AND ADD THESE LINES ***
			cashFlowStatements.POST("", handlers.CreateCashFlowStatement)      // POST /api/v1/cash-flow-statements
			cashFlowStatements.GET("", handlers.GetCashFlowStatements)         // GET /api/v1/cash-flow-statements
			cashFlowStatements.GET(":id", handlers.GetCashFlowStatement)       // GET /api/v1/cash-flow-statements/:id
			cashFlowStatements.PUT(":id", handlers.UpdateCashFlowStatement)    // PUT /api/v1/cash-flow-statements/:id
			cashFlowStatements.DELETE(":id", handlers.DeleteCashFlowStatement) // DELETE /api/v1/cash-flow-statements/:id
		}

		// --- Beneficial Owner Routes ---
		beneficialOwners := v1.Group("/beneficial-owners")
		{
			// *** UNCOMMENT AND ADD THESE LINES ***
			beneficialOwners.POST("", handlers.CreateBeneficialOwner)      // POST /api/v1/beneficial-owners
			beneficialOwners.GET("", handlers.GetBeneficialOwners)         // GET /api/v1/beneficial-owners
			beneficialOwners.GET(":id", handlers.GetBeneficialOwner)       // GET /api/v1/beneficial-owners/:id
			beneficialOwners.PUT(":id", handlers.UpdateBeneficialOwner)    // PUT /api/v1/beneficial-owners/:id
			beneficialOwners.DELETE(":id", handlers.DeleteBeneficialOwner) // DELETE /api/v1/beneficial-owners/:id
		}

		// --- Balance Sheet Routes ---
		balanceSheets := v1.Group("/balance-sheets")
		{
			// *** UNCOMMENT AND ADD THESE LINES ***
			balanceSheets.POST("", handlers.CreateBalanceSheet)      // POST /api/v1/balance-sheets
			balanceSheets.GET("", handlers.GetBalanceSheets)         // GET /api/v1/balance-sheets
			balanceSheets.GET(":id", handlers.GetBalanceSheet)       // GET /api/v1/balance-sheets/:id
			balanceSheets.PUT(":id", handlers.UpdateBalanceSheet)    // PUT /api/v1/balance-sheets/:id
			balanceSheets.DELETE(":id", handlers.DeleteBalanceSheet) // DELETE /api/v1/balance-sheets/:id
		}

		v1 := router.Group("/api/v1")
		{
			// ... (существующие маршруты для register, members, etc.) ...

			// Добавляем маршрут для детального поиска
			search := v1.Group("/search")
			{
				search.GET("/detailed", handlers.DetailedSearch) // GET /api/v1/search/detailed?q=...
				// Можно оставить и предыдущий поиск, если нужно
				// search.GET("/companies", handlers.SearchCompanies)
			}
		}
	}

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

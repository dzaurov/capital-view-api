package db

import (
	"capital-view-api/models" // Adjust import path if needed
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// Use the DSN from your Prisma schema
	dsn := "mydata.db"
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// AutoMigrate the schema
	// This creates the tables based on your Go structs if they don't exist
	err = DB.AutoMigrate(
		&models.Register{},
		&models.Member{},
		&models.IncomeStatement{},
		&models.FinancialStatement{},
		&models.CashFlowStatement{},
		&models.BeneficialOwner{},
		&models.BalanceSheet{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrated")
}

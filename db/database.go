// db/database.go
package db

import (
	"log" // Убедитесь, что log импортирован
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger" // Закомментировано в вашем коде, но может понадобиться
)

var DB *gorm.DB

func ConnectDatabase() error {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "mydata.db" // Имя по умолчанию, если DB_PATH не задан
	}

	// --- ДОБАВЬТЕ ЭТОТ ЛОГ ---
	// Чтобы точно знать, какой путь используется
	log.Printf("INFO: Attempting to connect to database at path: [%s]", dbPath)
	// -------------------------

	// Для теста можно временно закомментировать строки выше и использовать АБСОЛЮТНЫЙ ПУТЬ:
	// dbPath = "/Users/zeropera/Documents/capital-view-api/mydata.db"
	// log.Printf("INFO: [TEST] Using absolute database path: [%s]", dbPath)

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info), // Можно раскомментировать для GORM логов
	})

	if err != nil {
		log.Printf("ERROR: Failed to connect to database at path [%s]: %v", dbPath, err) // Добавим путь в лог ошибки
		return err
	}

	DB = database
	// log.Println("Database connection established.") // Можно раскомментировать
	return nil
}

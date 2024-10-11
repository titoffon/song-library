package migrations

import (
	"log"
	"os"

	"gorm.io/gorm"
)

// ExecuteSQLFile читает SQL-файл и выполняет запросы
func ExecuteSQLFile(db *gorm.DB, filepath string) error {

	query, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Выполнение SQL-запроса
	if err := db.Exec(string(query)).Error; err != nil {
		return err
	}

	log.Printf("SQL запрос из файла %s успешно выполнен", filepath)
	return nil
}

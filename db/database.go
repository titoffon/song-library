package db

import (
	"fmt"
	"log"
	"main/migrations"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB — глобальная переменная для работы с базой данных
var DB *gorm.DB

// InitDB инициализирует подключение к базе данных
func InitDB() {
	// Получение параметров подключения из переменных окружения
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}
	log.Println("Подключение к базе данных успешно")

	/*Ниже реализовано создание и заполнение тестовыми данными
	таблицы songs, которая хранит в себе всю необходимую информацию о песне.
	По завершениею работы программытестовые данные из таблицы удаляются */

	// Выполняем SQL-файл для создания таблицы
	if err := migrations.ExecuteSQLFile(DB, "migrations/create_songs_table.sql"); err != nil {
		log.Fatalf("Ошибка при выполнении SQL-файла: %v", err)
	}

	// Выполняем SQL-файл для заполнения таблицы
	if err := migrations.ExecuteSQLFile(DB, "migrations/insert_to_songs_table.sql"); err != nil {
		log.Fatalf("Ошибка при выполнении SQL-файла: %v", err)
	}

}

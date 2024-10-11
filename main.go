package main

import (
	"io"
	"log"
	"main/db"
	"net/http"
	"os"

	_ "main/docs" //импорт для документов Swagger

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Song Library API
// @version 1.0
// @description API для управления библиотекой песен
// @host localhost:8000
// @BasePath /songs
func main() {

	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Открываем файл для логирования
	file, err := os.OpenFile("debug/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для логов: %v", err)
	}
	defer file.Close()

	// Настраиваем MultiWriter для записи и в файл, и в консоль
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	// Инициализация базы данных
	db.InitDB()

	// Настройка маршрутизатора
	r := SetupRouter()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000" // Значение по умолчанию
	}

	// Регистрируем обработчик Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

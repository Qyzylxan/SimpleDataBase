package main

import (
	"log"
	"net/http"
	"os"

	"simpledatabase/internal/database"
	"simpledatabase/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используются переменные окружения системы")
	}

	// Инициализация базы данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}
	defer db.Close()

	// Инициализация обработчиков
	handlers.InitHandlers(db)

	// Настройка маршрутов
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/submit", handlers.SubmitHandler)
	http.HandleFunc("/data", handlers.DataHandler)
	http.HandleFunc("/delete/", handlers.DeleteHandler) // Новый маршрут для удаления
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

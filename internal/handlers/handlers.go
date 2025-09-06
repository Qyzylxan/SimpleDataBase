package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"simpledatabase/internal/database"
	"simpledatabase/internal/models"
	"strconv"
	"strings"
)

var db *sql.DB
var templates *template.Template

// InitHandlers инициализирует обработчики
func InitHandlers(database *sql.DB) {
	db = database
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

// IndexHandler обрабатывает главную страницу
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// SubmitHandler обрабатывает отправку формы
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Парсинг формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		// Валидация данных
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		if firstName == "" || lastName == "" || email == "" || phone == "" {
			http.Error(w, "Все поля обязательны для заполнения", http.StatusBadRequest)
			return
		}

		// Вставка данных в базу
		id, err := database.InsertUser(firstName, lastName, email, phone)
		if err != nil {
			log.Printf("Ошибка вставки пользователя: %v", err)
			http.Error(w, "Ошибка сохранения данных", http.StatusInternalServerError)
			return
		}

		// Ответ в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Данные успешно сохранены",
			"id":      id,
		})
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// DataHandler обрабатывает страницу с данными
func DataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Получение данных из базы
		users, err := database.GetAllUsers()
		if err != nil {
			log.Printf("Ошибка получения пользователей: %v", err)
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}

		// Отображение шаблона с данными
		err = templates.ExecuteTemplate(w, "data.html", users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// DeleteHandler обрабатывает удаление пользователей
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		// Обработка массового удаления через JSON
		if r.Header.Get("Content-Type") == "application/json" {
			var deleteReq models.DeleteRequest
			err := json.NewDecoder(r.Body).Decode(&deleteReq)
			if err != nil {
				sendJSONResponse(w, false, "Неверный формат JSON", http.StatusBadRequest)
				return
			}

			if len(deleteReq.UserIDs) == 0 {
				sendJSONResponse(w, false, "Не указаны ID пользователей для удаления", http.StatusBadRequest)
				return
			}

			// Удаление пользователей из базы
			rowsAffected, err := database.DeleteUsers(deleteReq.UserIDs)
			if err != nil {
				sendJSONResponse(w, false, "Ошибка удаления пользователей: "+err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(models.DeleteResponse{
				Success:      true,
				Message:      "Пользователи успешно удалены",
				DeletedCount: int(rowsAffected),
			})
			return
		}

		// Обработка одиночного удаления через URL (для обратной совместимости)
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) >= 3 {
			idStr := pathParts[2]
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				sendJSONResponse(w, false, "Неверный ID пользователя", http.StatusBadRequest)
				return
			}

			rowsAffected, err := database.DeleteUser(id)
			if err != nil {
				sendJSONResponse(w, false, "Ошибка удаления пользователя: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if rowsAffected == 0 {
				sendJSONResponse(w, false, "Пользователь не найден", http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(models.DeleteResponse{
				Success:      true,
				Message:      "Пользователь успешно удалён",
				DeletedCount: 1,
			})
			return
		}

		sendJSONResponse(w, false, "Неверный запрос", http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, false, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// Вспомогательная функция для отправки JSON ответов
func sendJSONResponse(w http.ResponseWriter, success bool, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: success,
		Message: message,
	})
}

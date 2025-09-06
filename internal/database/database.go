package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"simpledatabase/internal/models"
	"strings"

	_ "github.com/lib/pq"
)

// DB глобальная переменная базы данных
var DB *sql.DB

// InitDB инициализирует подключение к базе данных и создает таблицы
func InitDB() (*sql.DB, error) {
	// Подключение к PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Проверка подключения
	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %v", err)
	}

	// Создание таблицы если не существует
	err = createTables()
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблиц: %v", err)
	}

	log.Println("База данных подключена успешно")
	return DB, nil
}

// createTables создает необходимые таблицы в базе данных
func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		phone VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("ошибка выполнения SQL: %v", err)
	}

	return nil
}

// InsertUser вставляет нового пользователя в базу данных
func InsertUser(firstName, lastName, email, phone string) (string, error) {
	insertSQL := `
	INSERT INTO users (first_name, last_name, email, phone)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (email) DO UPDATE SET
	first_name = EXCLUDED.first_name,
	last_name = EXCLUDED.last_name,
	phone = EXCLUDED.phone
	RETURNING id
	`

	var id string
	err := DB.QueryRow(insertSQL, firstName, lastName, email, phone).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("ошибка вставки пользователя: %v", err)
	}

	return id, nil
}

// GetAllUsers возвращает всех пользователей из базы данных
func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`
		SELECT id, first_name, last_name, email, phone 
		FROM users 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %v", err)
	}

	return users, nil
}

// DeleteUser удаляет пользователя по ID
func DeleteUser(id int) (int64, error) {
	result, err := DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("ошибка удаления пользователя: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества удаленных строк: %v", err)
	}

	return rowsAffected, nil
}

// DeleteUsers удаляет нескольких пользователей по массиву ID
func DeleteUsers(userIDs []int) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	// Перевод массива ID в строку для формирования запроса
	idString := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs)) // Странный массив args
	for i, id := range userIDs {
		idString[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// Формирование SQL-запроса (данные string в одинарных кавычках!)
	query := fmt.Sprintf("DELETE FROM users WHERE id IN (%s)", strings.Join(idString, ","))

	result, err := DB.Exec(query, args...) // Который тут участвует в подстановке в id?
	if err != nil {
		return 0, fmt.Errorf("ошибка удаления пользователей: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества удаленных строк: %v", err)
	}

	return rowsAffected, nil
}

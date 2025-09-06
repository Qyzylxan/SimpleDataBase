package database

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	// Тест можно запустить только при наличии тестовой базы данных
	t.Skip("Требуется тестовая база данных")

	_, err := InitDB()
	if err != nil {
		t.Errorf("InitDB() error = %v", err)
	}
}

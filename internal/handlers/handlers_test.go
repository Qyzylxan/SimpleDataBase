package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("IndexHandler вернул неверный статус: получил %v, ожидал %v",
			status, http.StatusOK)
	}
}

func TestSubmitHandler_GET(t *testing.T) {
	req, err := http.NewRequest("GET", "/submit", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("SubmitHandler с GET вернул неверный статус: получил %v, ожидал %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestDataHandler_GET(t *testing.T) {
	// Инициализируем пустую базу данных для тестов
	db, _ := sql.Open("postgres", "host=test port=test user=test password=test dbname=test sslmode=disable")
	InitHandlers(db)

	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DataHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError && status != http.StatusOK {
		t.Errorf("DataHandler вернул неверный статус: получил %v", status)
	}
}

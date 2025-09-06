package routers

import (
	"net/http"
	"simpledatabase/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	router.HandleFunc("/submit", handlers.SubmitHandler).Methods("POST")
	router.HandleFunc("/data", handlers.DataHandler).Methods("GET")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
}

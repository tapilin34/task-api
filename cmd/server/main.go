package main

import (
	"log"
	"net/http"
	"tasks-api/internal/handlers"
	"tasks-api/internal/storage"
)

func main() {
	// Подключаем конкретную реализацию (in‑memory) интерфейса Storage
	var store storage.Storage
	store = storage.NewCache()
	h := handlers.New(store)
	// добавляем новые эндпоинты
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE

	log.Println("server listening on http://localhost:8081/tasks")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}

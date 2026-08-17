package main

import (
	"log"
	"net/http"
	"tasks-api/internal/handlers"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
	"time"
)

func main() {
	// Подключаем конкретную реализацию (in‑memory) интерфейса Storage
	var store storage.Storage
	store = storage.NewCache()
	store.Create(models.Task{ID: 1, Title: "Buy tea", Done: true, CreatedAt: time.Now().Format(time.RFC3339)})
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

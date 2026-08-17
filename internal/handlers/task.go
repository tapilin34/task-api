package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }
type ErrorResponse struct {
	Error string
}

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("[TasksCollection] start, method:%s", r.Method)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		log.Print("[TasksCollection] GET")
		list := h.Store.List()
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		log.Print("[TasksCollection] POST: start decode")

		var m models.Task
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			log.Printf("[TasksCollection] POST: decode error:%s", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат данных"})
			return
		}

		log.Printf("[TasksCollection] POST: decoded task id:%d", m.ID)

		if m.Title == "" {
			log.Print("[TasksCollection] POST: validation error (empty title)")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Вы не заполнили заголовок task"})
			return
		}

		log.Print("[TasksCollection] POST: calling Store.Create")
		task, err := h.Store.Create(m)
		if err != nil {
			log.Printf("[TasksCollection] POST: Store.Create error:%s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("[TasksCollection] POST: created task: %s", task.Title)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
		log.Print("[TasksCollection] POST: finished")

	default:
		log.Printf("[TasksCollection] unsupported method:%s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Метод не поддерживается"})
	}

	log.Print("[TasksCollection] end")
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("[TaskItem] start, method:%s, path %s", r.Method, r.URL.Path)
	// парсим url
	urls := strings.Split(r.URL.Path, "/")
	if len(urls) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный путь"})
		return
	}
	// приведение к int
	id, _ := strconv.Atoi(urls[2])
	// логика для разных методов
	switch r.Method {
	case http.MethodGet:
		list, ok := h.Store.Get(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPut:
		var m models.Task
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат данных"})
			return
		}
		upd, err := h.Store.Update(id, m)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(upd)
	case http.MethodDelete:
		err := h.Store.Delete(id)
		if err != nil {
			if err == storage.ErrTaskNotFound {
				w.WriteHeader(http.StatusNotFound)
			}
		}
		w.WriteHeader(http.StatusAccepted)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Метод не поддерживается"})
	}
}

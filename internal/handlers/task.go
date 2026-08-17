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
	Error string `json:"error"`
}

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Printf("[TasksCollection] start, method:%s, path:%s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		log.Print("[TasksCollection] GET")
		list := h.Store.List()
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		if !isJSONRequest(r) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Content-Type должен быть application/json"})
			return
		}
		log.Print("[TasksCollection] POST: start decode")

		var m models.Task
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			log.Printf("[TasksCollection] POST: decode error:%s", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат данных"})
			return
		}

		log.Printf("[TasksCollection] POST: decoded task id:%d", m.ID)

		if strings.TrimSpace(m.Title) == "" {
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Printf("[TaskItem] start, method:%s, path %s", r.Method, r.URL.Path)
	// парсим url
	urls := strings.Split(r.URL.Path, "/")
	if len(urls) != 3 || urls[1] != "tasks" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный путь"})
		return
	}
	// приведение к int
	id, idErr := strconv.Atoi(urls[2])
	if idErr != nil {
		// не число → 400
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат ID"})
		return
	}
	// если ID начинаются с 1, то 0 и отрицательные считаем некорректными
	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	// логика для разных методов
	switch r.Method {
	case http.MethodGet:
		list, ok := h.Store.Get(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
			return
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPut:
		if !isJSONRequest(r) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Content-Type должен быть application/json"})
			return
		}
		var m models.Task
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат данных"})
			return
		}
		if strings.TrimSpace(m.Title) == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Вы не заполнили заголовок task"})
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
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
				return
			}
			// любая другая ошибка → 500, как в POST/PUT
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		// успешное удаление → 204 No Content без тела
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Метод не поддерживается"})
	}
}

func isJSONRequest(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	// допускаем "application/json" и варианты с charset и т.п.
	return strings.HasPrefix(ct, "application/json")
}

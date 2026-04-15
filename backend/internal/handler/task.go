package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"

	"github.com/Ruturaj-7802/taskflow/internal/dto"
	"github.com/Ruturaj-7802/taskflow/internal/middleware"
	"github.com/Ruturaj-7802/taskflow/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// helper to extract userID safely
func getUserID(r *http.Request) (uuid.UUID, error) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		return uuid.Nil, http.ErrNoCookie
	}
	return uuid.Parse(userIDStr)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	ownerID, err := getUserID(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	t, err := h.svc.Create(r.Context(), projectID, req, ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	ownerID, err := getUserID(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	tasks, err := h.svc.List(r.Context(), projectID, ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	ownerID, err := getUserID(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	t, err := h.svc.Get(r.Context(), id, ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		if err.Error() == "not found" {
			status = http.StatusNotFound
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	ownerID, err := getUserID(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	t, err := h.svc.Update(r.Context(), id, req, ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		if err.Error() == "not found" {
			status = http.StatusNotFound
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ownerID, err := getUserID(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id, ownerID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		if err.Error() == "not found" {
			status = http.StatusNotFound
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
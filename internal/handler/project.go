package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/ruturaj/taskflow/internal/dto"
    "github.com/ruturaj/taskflow/internal/middleware"
    "github.com/ruturaj/taskflow/internal/service"
)

type ProjectHandler struct {
    projectSvc *service.ProjectService
}

func NewProjectHandler(projectSvc *service.ProjectService) *ProjectHandler {
    return &ProjectHandler{projectSvc: projectSvc}
}

func currentUserID(r *http.Request) uuid.UUID {
    claims := r.Context().Value(middleware.UserClaimsKey).(*service.Claims)
    return claims.UserID
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
    projects, err := h.projectSvc.List(r.Context(), currentUserID(r))
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"projects": projects})
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateProjectRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    if req.Name == "" {
        writeValidationError(w, map[string]string{"name": "is required"})
        return
    }

    p, err := h.projectSvc.Create(r.Context(), currentUserID(r), req)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    writeJSON(w, http.StatusCreated, p)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid project id")
        return
    }

    p, err := h.projectSvc.Get(r.Context(), id, currentUserID(r))
    if err != nil {
        if errors.Is(err, service.ErrNotFound) {
            writeError(w, http.StatusNotFound, "not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    writeJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid project id")
        return
    }

    var req dto.UpdateProjectRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    p, err := h.projectSvc.Update(r.Context(), id, currentUserID(r), req)
    if err != nil {
        if errors.Is(err, service.ErrNotFound) {
            writeError(w, http.StatusNotFound, "not found")
            return
        }
        if errors.Is(err, service.ErrForbidden) {
            writeError(w, http.StatusForbidden, "forbidden")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    writeJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid project id")
        return
    }

    err = h.projectSvc.Delete(r.Context(), id, currentUserID(r))
    if err != nil {
        if errors.Is(err, service.ErrNotFound) {
            writeError(w, http.StatusNotFound, "not found")
            return
        }
        if errors.Is(err, service.ErrForbidden) {
            writeError(w, http.StatusForbidden, "forbidden")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) Stats(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid project id")
        return
    }

    stats, err := h.projectSvc.Stats(r.Context(), id, currentUserID(r))
    if err != nil {
        if errors.Is(err, service.ErrNotFound) {
            writeError(w, http.StatusNotFound, "not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    writeJSON(w, http.StatusOK, stats)
}
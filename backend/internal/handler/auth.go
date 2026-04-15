package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/Ruturaj-7802/taskflow/internal/dto"
    "github.com/Ruturaj-7802/taskflow/internal/service"
)

type AuthHandler struct {
    authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
    return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req dto.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // Validation
    fields := map[string]string{}
    if req.Name == "" { fields["name"] = "is required" }
    if req.Email == "" { fields["email"] = "is required" }
    if len(req.Password) < 8 { fields["password"] = "must be at least 8 characters" }
    if len(fields) > 0 {
        writeValidationError(w, fields)
        return
    }

    resp, err := h.authSvc.Register(r.Context(), req)
    if err != nil {
        if errors.Is(err, service.ErrEmailTaken) {
            writeValidationError(w, map[string]string{"email": "already registered"})
            return
        }
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }

    writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    resp, err := h.authSvc.Login(r.Context(), req)
    if err != nil {
        if errors.Is(err, service.ErrInvalidCreds) {
            writeError(w, http.StatusUnauthorized, "invalid email or password")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }

    writeJSON(w, http.StatusOK, resp)
}
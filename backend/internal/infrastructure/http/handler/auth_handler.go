package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	appuser "github.com/praxisvr/grama/internal/application/user"
	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/infrastructure/auth"
	"github.com/praxisvr/grama/internal/infrastructure/http/middleware"
	"github.com/praxisvr/grama/internal/ports/repository"
)

// AuthHandler handles all authentication endpoints.
type AuthHandler struct {
	registerUC *appuser.RegisterUseCase
	loginUC    *appuser.LoginUseCase
	refreshUC  *appuser.RefreshUseCase
	tokenRepo  repository.TokenRepository
}

func NewAuthHandler(
	registerUC *appuser.RegisterUseCase,
	loginUC *appuser.LoginUseCase,
	refreshUC *appuser.RefreshUseCase,
	tokenRepo repository.TokenRepository,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
		tokenRepo:  tokenRepo,
	}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		Role         string `json:"role"`
		ConsentGiven bool   `json:"consent_given"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	out, err := h.registerUC.Execute(r.Context(), appuser.RegisterInput{
		Name:         req.Name,
		Email:        req.Email,
		Password:     req.Password,
		Role:         req.Role,
		ConsentGiven: req.ConsentGiven,
	})
	if err != nil {
		switch {
		case errors.Is(err, appuser.ErrConsentRequired):
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "consent required"})
		case errors.Is(err, appuser.ErrOperatorSelfRegister):
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "operators cannot self-register"})
		case errors.Is(err, domainuser.ErrEmailInvalid),
			errors.Is(err, domainuser.ErrPasswordTooShort),
			errors.Is(err, domainuser.ErrInvalidRole):
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "registration failed"})
		}
		return
	}

	// NEVER include password_hash in any HTTP response.
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         out.ID,
		"name":       out.Name,
		"email":      out.Email,
		"role":       out.Role,
		"created_at": out.CreatedAt,
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	out, err := h.loginUC.Execute(r.Context(), appuser.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, appuser.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "login failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  out.AccessToken,
		"refresh_token": out.RefreshToken,
		"expires_in":    out.ExpiresIn,
		"user": map[string]string{
			"id":   out.UserID,
			"name": out.UserName,
			"role": out.UserRole,
		},
	})
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	out, err := h.refreshUC.Execute(r.Context(), appuser.RefreshInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, appuser.ErrInvalidToken) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "refresh failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  out.AccessToken,
		"refresh_token": out.RefreshToken,
		"expires_in":    out.ExpiresIn,
	})
}

// Logout handles POST /auth/logout. Requires a valid JWT (applied by AuthMiddleware).
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if _, ok := middleware.ClaimsFromContext(r.Context()); !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	hash := auth.HashRefreshToken(req.RefreshToken)
	if err := h.tokenRepo.RevokeRefreshToken(r.Context(), hash); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "logout failed"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writeJSON encode: %v", err)
	}
}

package handler

import (
	"errors"
	"net/http"

	appuser "github.com/praxisvr/grama/internal/application/user"
	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/infrastructure/auth"
)

// UserHandler handles profile and operator management endpoints.
type UserHandler struct {
	getProfileUC    *appuser.GetProfileUseCase
	updateProfileUC *appuser.UpdateProfileUseCase
	createOpUC      *appuser.CreateOperatorUseCase
	listOpsUC       *appuser.ListOperatorsUseCase
}

func NewUserHandler(
	getProfileUC *appuser.GetProfileUseCase,
	updateProfileUC *appuser.UpdateProfileUseCase,
	createOpUC *appuser.CreateOperatorUseCase,
	listOpsUC *appuser.ListOperatorsUseCase,
) *UserHandler {
	return &UserHandler{
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
		createOpUC:      createOpUC,
		listOpsUC:       listOpsUC,
	}
}

// GetMe handles GET /users/me.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	out, err := h.getProfileUC.Execute(r.Context(), claims.Subject)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch profile"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":         out.ID,
		"name":       out.Name,
		"email":      out.Email,
		"role":       out.Role,
		"created_at": out.CreatedAt,
	})
}

// UpdateMe handles PATCH /users/me.
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.updateProfileUC.Execute(r.Context(), claims.Subject, appuser.UpdateProfileInput{
		Name: req.Name,
	}); err != nil {
		if errors.Is(err, appuser.ErrNameRequired) {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "name is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "profile updated"})
}

// CreateOperator handles POST /users/operators (owner only, enforced by RequireRole middleware).
func (h *UserHandler) CreateOperator(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	out, err := h.createOpUC.Execute(r.Context(), claims.Subject, appuser.CreateOperatorInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrEmailInvalid),
			errors.Is(err, domainuser.ErrPasswordTooShort):
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "operator creation failed"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":    out.ID,
		"name":  out.Name,
		"email": out.Email,
		"role":  out.Role,
	})
}

// ListOperators handles GET /users/operators (owner only, enforced by RequireRole middleware).
func (h *UserHandler) ListOperators(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	out, err := h.listOpsUC.Execute(r.Context(), claims.Subject)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not list operators"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"operators": out.Operators})
}

package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"app/domain"
	"app/internal/auth"
)

type Handler struct {
	Service *Service
	JWT     *auth.JWT
}

func NewHandler(s *Service, jwt *auth.JWT) *Handler {
	return &Handler{Service: s, JWT: jwt}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var dto domain.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if dto.Email == "" || dto.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	u, err := h.Service.Register(
		r.Context(),
		dto.Email,
		dto.Password,
		dto.Username,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.JWT.Generate(u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	resp := domain.AuthResponse{
		User:  ToUserResponse(u),
		Token: token,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var dto domain.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if dto.Email == "" || dto.Password == "" {
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	u, err := h.Service.Login(r.Context(), dto.Email, dto.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			respondWithError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	token, err := h.JWT.Generate(u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	resp := domain.AuthResponse{
		User:  ToUserResponse(u),
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"app/domain"
	"app/internal/auth"
	"app/internal/refreshtoken"
)

type Handler struct {
	Service        *Service
	JWT            *auth.JWT
	RefreshService *refreshtoken.Service
}

func NewHandler(s *Service, jwt *auth.JWT, refreshService *refreshtoken.Service) *Handler {
	return &Handler{
		Service:        s,
		JWT:            jwt,
		RefreshService: refreshService,
	}
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

	accessToken, err := h.JWT.Generate(u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refreshToken, err := h.RefreshService.Generate(r.Context(), u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	resp := domain.AuthResponse{
		User:         ToUserResponse(u),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	accessToken, err := h.JWT.Generate(u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refreshToken, err := h.RefreshService.Generate(r.Context(), u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	resp := domain.AuthResponse{
		User:         ToUserResponse(u),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var dto domain.RefreshTokenDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if dto.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	newRefreshToken, userID, err := h.RefreshService.Rotate(r.Context(), dto.RefreshToken)
	if err != nil {
		if errors.Is(err, refreshtoken.ErrInvalidRefreshToken) ||
			errors.Is(err, refreshtoken.ErrExpiredRefreshToken) ||
			errors.Is(err, refreshtoken.ErrRevokedRefreshToken) {
			respondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	u, err := h.Service.GetByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	accessToken, err := h.JWT.Generate(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	resp := domain.AuthResponse{
		User:         ToUserResponse(u),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var dto domain.LogoutDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if dto.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	if err := h.RefreshService.Revoke(r.Context(), dto.RefreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

package user

import (
	"encoding/json"
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
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	u, err := h.Service.Register(
		r.Context(),
		dto.Email,
		dto.Password,
		dto.Username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, _ := h.JWT.Generate(u.ID)

	resp := struct {
		User  domain.UserResponse `json:"user"`
		Token string              `json:"token"`
	}{
		User:  ToUserResponse(u),
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

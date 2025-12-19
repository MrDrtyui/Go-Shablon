package router

import (
	"app/internal/user"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi *chi.Mux
}

func NewRouter(userHandler *user.Handler) *Router {
	r := chi.NewRouter()

	r.Post("/auth/register", userHandler.Register)

	return &Router{chi: r}
}

func (r *Router) Handler() http.Handler {
	return r.chi
}

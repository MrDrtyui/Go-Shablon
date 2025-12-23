package router

import (
	"app/internal/user"
	"net/http"

	_ "app/docs" // Import generated docs

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	chi *chi.Mux
}

func NewRouter(userHandler *user.Handler) *Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
		r.Post("/logout", userHandler.Logout)
	})

	return &Router{chi: r}
}

func (r *Router) Handler() http.Handler {
	return r.chi
}

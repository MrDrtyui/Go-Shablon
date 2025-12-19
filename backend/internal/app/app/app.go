package app

import (
	"net/http"

	"app/internal/auth"
	"app/internal/config"
	"app/internal/db"
	"app/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/labstack/gommon/log"
)

func Run() error {
	// TODO: init config
	cfg := config.MustLoad()
	log.Info(cfg.Env)

	// TODO: init db
	db := db.NewDbClient(cfg.Database.Url)
	log.Info("Open connect to db", db.Client.Schema)
	defer db.Close()

	// TODO: jwt
	jwtSvc := auth.New(cfg.Jwt.Secret, cfg.Jwt.TtlHours)

	// TODO: init chi
	userRepo := user.NewPostgresRepo(db)
	userService := user.NewSercie(userRepo)
	userHandler := user.NewHandler(userService, jwtSvc)

	chi := chi.NewRouter()

	chi.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello epta"))
	})

	chi.Post("/auth/register", userHandler.Register)

	// TODO: run server
	log.Info("Starting server at", cfg.Port)
	return http.ListenAndServe(cfg.Port, chi)
}

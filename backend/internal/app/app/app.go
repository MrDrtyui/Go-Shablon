package app

import (
	"app/internal/auth"
	"app/internal/config"
	"app/internal/db"
	"app/internal/router"
	"app/internal/user"

	"github.com/labstack/gommon/log"
)

type App struct {
	Router *router.Router
	Cfg    *config.Config
	Db     *db.Db
}

func NewApp() *App {
	// TODO: init config
	cfg := config.MustLoad()
	log.Info(cfg.Env)

	// TODO: init db
	db := db.NewDbClient(cfg.Database.Url)
	log.Info("Open connect to db", db.Client.Schema)

	// TODO: jwt
	jwtSvc := auth.New(cfg.Jwt.Secret, cfg.Jwt.TtlHours)

	// TODO: init chi
	userRepo := user.NewPostgresRepo(db)
	userService := user.NewSercie(userRepo)
	userHandler := user.NewHandler(userService, jwtSvc)

	r := router.NewRouter(userHandler)

	return &App{
		Router: r,
		Cfg:    cfg,
		Db:     db,
	}
}

package app

import (
	"net/http"

	"github.com/labstack/gommon/log"
)

func Run() {
	app := NewApp()
	defer app.Db.Close()

	log.Info("Start server on ", app.Cfg.Http.Port)

	if err := http.ListenAndServe(app.Cfg.Http.Port, app.Router.Handler()); err != nil {
		log.Fatal("Failed to run http server", err.Error())
	}
}

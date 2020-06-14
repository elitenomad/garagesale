package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// App is the entrypoint to web application
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp knows how to construct internal state for an App.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}

func (app *App) Handle(method, pattern string, fn http.HandlerFunc) {
	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

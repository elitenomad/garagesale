package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler func(http.ResponseWriter, *http.Request) error

// App is the entrypoint to web application
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
}

// NewApp knows how to construct internal state for an App.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
		mw:  mw,
	}
}

func (app *App) Handle(method, pattern string, h Handler) {

	h = wrapMiddleware(app.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		if err != nil {
			app.log.Printf("Unhandled ERROR : %+v", err)
		}
	}

	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

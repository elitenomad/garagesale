package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)


type Handler func(http.ResponseWriter, *http.Request) error

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

func (app *App) Handle(method, pattern string, h Handler) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		if err != nil {
			app.log.Printf("ERROR : %+v", err)

			if err := RespondError(w, err); err != nil {
				app.log.Printf("ERROR : %v", err)
			}
		}
	}

	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

type Handler func(context.Context, http.ResponseWriter, *http.Request) error

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

		v := Values{
			Start: time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		err := h(ctx, w, r)

		if err != nil {
			app.log.Printf("Unhandled ERROR : %+v", err)
		}
	}

	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

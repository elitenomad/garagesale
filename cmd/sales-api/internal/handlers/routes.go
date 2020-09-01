package handlers

import (
	"log"
	"net/http"

	mid "github.com/elitenomad/garagesale/internal/middleware"
	"github.com/elitenomad/garagesale/internal/platform/auth"
	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

func API(logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(logger, mid.Errors(logger), mid.Metrics())

	{
		c := Check{
			Db:  db,
			Log: logger,
		}

		app.Handle(http.MethodGet, "/api/v1/health", c.Health)
	}

	{
		u := Users{
			DB:            db,
			authenticator: authenticator,
		}

		app.Handle(http.MethodGet, "/api/v1/users/token", u.Token)
	}

	{
		p := Product{
			Db:  db,
			Log: logger,
		}

		app.Handle(http.MethodGet, "/api/v1/products", p.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/api/v1/products/{id}", p.Fetch, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/api/v1/products", p.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/api/v1/products/{id}", p.Update, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/api/v1/products/{id}", p.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

		app.Handle(http.MethodPost, "/api/v1/products/{id}/sales", p.AddSale, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/api/v1/products/{id}/sales", p.ListSales, mid.Authenticate(authenticator))
	}

	return app
}

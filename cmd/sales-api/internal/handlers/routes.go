package handlers

import (
	"log"
	"net/http"

	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)

	p := Product{
		Db:  db,
		Log: logger,
	}

	c := Check{
		Db:  db,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/api/v1/health", c.Health)
	app.Handle(http.MethodGet, "/api/v1/products", p.List)
	app.Handle(http.MethodGet, "/api/v1/products/{id}", p.Fetch)
	app.Handle(http.MethodPost, "/api/v1/products", p.Create)
	app.Handle(http.MethodPut, "/api/v1/products/{id}", p.Update)
	app.Handle(http.MethodDelete, "/api/v1/products/{id}", p.Delete)

	app.Handle(http.MethodPost, "/api/v1/products/{id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/api/v1/products/{id}/sales", p.ListSales)

	return app
}

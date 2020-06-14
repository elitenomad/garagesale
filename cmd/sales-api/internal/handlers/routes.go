package handlers

import (
	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)

	p := Product{
		Db: db,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/api/v1/products", p.List)
	app.Handle(http.MethodGet, "/api/v1/products/{id}", p.Fetch)

	return app
}

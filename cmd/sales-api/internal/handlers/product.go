package handlers

import (
	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/elitenomad/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type Products struct {
	Db *sqlx.DB
	Log *log.Logger
}

type Product struct {
	Db *sqlx.DB
	Log *log.Logger
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) error {
	products, err := product.List(p.Db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(w, products, http.StatusOK)
}

func (p *Product) Fetch(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	pdct, err := product.Fetch(p.Db, id)
	if err != nil {
		if err != nil {
			switch err {
			case product.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case product.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "getting product %q", id)
			}
		}
	}

	return web.Respond(w, pdct, http.StatusOK)
}


func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	product, err := product.Create(p.Db, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(w, product, http.StatusCreated)
}

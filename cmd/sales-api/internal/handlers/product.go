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
	products, err := product.List(r.Context(), p.Db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(w, products, http.StatusOK)
}

func (p *Product) Fetch(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	pdct, err := product.Fetch(r.Context(), p.Db, id)
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

	product, err := product.Create(r.Context(), p.Db, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(w, product, http.StatusCreated)
}

func (p *Product) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(r.Context(), p.Db, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(w, nil, http.StatusNoContent)
}

func (p *Product) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := product.Delete(r.Context(), p.Db, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(w, nil, http.StatusNoContent)
}

func (p *Product) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(r.Context(), p.Db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Product) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(r.Context(), p.Db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(w, list, http.StatusOK)
}

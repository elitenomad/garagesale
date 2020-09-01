package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/elitenomad/garagesale/internal/platform/auth"
	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/elitenomad/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type Product struct {
	Db  *sqlx.DB
	Log *log.Logger
}

func (p *Product) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(r.Context(), "internal.handlers.product.list")
	defer span.End()

	products, err := product.List(ctx, p.Db)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(w, products, http.StatusOK)
}

func (p *Product) Fetch(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	ctx, span := trace.StartSpan(r.Context(), "internal.handlers.product.fetch")
	defer span.End()

	pdct, err := product.Fetch(ctx, p.Db, id)
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

func (p *Product) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(r.Context(), "internal.handlers.product.Create")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)

	if !ok {
		return errors.New("Auth claims is not in context")
	}

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	product, err := product.Create(ctx, p.Db, claims, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(w, product, http.StatusCreated)
}

func (p *Product) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)

	if !ok {
		return errors.New("Auth claims is not in context")
	}

	if err := product.Update(ctx, p.Db, id, claims, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(w, nil, http.StatusNoContent)
}

func (p *Product) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(r.Context(), "internal.handlers.product.Delete")
	defer span.End()

	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.Db, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(w, nil, http.StatusNoContent)
}

func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(r.Context(), "internal.handlers.product.AddSale")
	defer span.End()

	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.Db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Product) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(r.Context(), "internal.handlers.product.ListSales")
	defer span.End()

	id := chi.URLParam(r, "id")

	list, err := product.ListSales(ctx, p.Db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(w, list, http.StatusOK)
}

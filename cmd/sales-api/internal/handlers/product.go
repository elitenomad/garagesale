package handlers

import (
	"github.com/elitenomad/garagesale/internal/platform/web"
	"github.com/elitenomad/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
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

func (p *Product) List(w http.ResponseWriter, r *http.Request) {
	products, err := product.List(p.Db)
	if err != nil {
		p.Log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Respond(w, products, http.StatusOK); err != nil {
		p.Log.Println("error writing result", err)
	}
}

func (p *Product) Fetch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := product.Fetch(p.Db, id)
	if err != nil {
		p.Log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Respond(w, product, http.StatusOK); err != nil {
		p.Log.Println("error writing result", err)
	}
}


func (p *Product) Create(w http.ResponseWriter, r *http.Request) {

	// Ensure we parse the info passed onto the API. In this case the Product data
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	product, err := product.Create(p.Db, np, time.Now())
	if err != nil {
		p.Log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Respond(w, product, http.StatusCreated); err != nil {
		p.Log.Println("error writing result", err)
	}
}

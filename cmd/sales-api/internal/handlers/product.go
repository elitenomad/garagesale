package handlers

import (
	"encoding/json"
	"github.com/elitenomad/garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
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

	data, err := json.Marshal(products)
	if err != nil {
		p.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
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

	data, err := json.Marshal(product)
	if err != nil {
		p.Log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing result", err)
	}
}

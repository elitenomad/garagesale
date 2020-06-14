package handlers

import (
	"encoding/json"
	"github.com/elitenomad/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Products struct {
	Db *sqlx.DB
}

func (p *Products) List(w http.ResponseWriter, r *http.Request) {
	products, err := product.List(p.Db)
	if err != nil {
		log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(products)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}

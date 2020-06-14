package product

import (
	"github.com/jmoiron/sqlx"
)

func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * from products`

	if err := db.Select(&products, q); err != nil {
		return nil, err
	}

	return products, nil
}

func Fetch(db *sqlx.DB, id string) (*Product, error) {
	var p Product

	// No inputs in the query
	const q = `SELECT * from products where product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		return nil, err
	}


	return &p, nil
}




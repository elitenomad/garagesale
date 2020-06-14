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


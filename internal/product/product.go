package product

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
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

func Create(db *sqlx.DB, data NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID: uuid.New().String(),
		Name:     data.Name,
		Cost:     data.Cost,
		Quantity: data.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	// No inputs in the query
	const q = `INSERT into products 
		(product_id, name, cost, quantity, date_created, date_updated) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(q, p.ID,p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "Insert product failed...")
	}

	return &p, nil
}




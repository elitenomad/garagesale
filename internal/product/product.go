package product

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

// Failure scenarios.
var (
	ErrNotFound = errors.New("product not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT * from products`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "Listing products")
	}

	return products, nil
}

func Fetch(db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Product

	// No inputs in the query
	const q = `SELECT * from products where product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single product")
	}


	return &p, nil
}

func Create(db *sqlx.DB, data NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID: uuid.New().String(),
		Name:     data.Name,
		Cost:     data.Cost,
		Quantity: data.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO products
		(product_id, name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(q,
		p.ID, p.Name,
		p.Cost, p.Quantity,
		p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}

	return &p, nil
}




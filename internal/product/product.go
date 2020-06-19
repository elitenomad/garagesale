package product

import (
	"context"
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

func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT
				p.*,
				COALESCE(SUM(s.quantity), 0) as sold,
				COALESCE(SUM(s.paid), 0) as revenue
			FROM products AS p
			LEFT JOIN sales AS s ON p.product_id = s.product_id
			GROUP BY p.product_id`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "Listing products")
	}

	return products, nil
}

func Fetch(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Product

	// No inputs in the query
	const q = `SELECT
				p.*,
				COALESCE(SUM(s.quantity), 0) as sold,
				COALESCE(SUM(s.paid), 0) as revenue
			FROM products AS p
			LEFT JOIN sales AS s ON p.product_id = s.product_id
			where p.product_id = $1
			GROUP BY p.product_id`

	if err := db.Get(&p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single product")
	}


	return &p, nil
}

func Create(ctx context.Context, db *sqlx.DB, data NewProduct, now time.Time) (*Product, error) {
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


func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateProduct, now time.Time) error {
	p, err := Fetch(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Cost != nil {
		p.Cost = *update.Cost
	}
	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}
	p.DateUpdated = now

	const q = `UPDATE products SET
		"name" = $2,
		"cost" = $3,
		"quantity" = $4,
		"date_updated" = $5
		WHERE product_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

func Delete(ctx context.Context, db *sqlx.DB, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM products WHERE product_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting product %s", id)
	}

	return nil
}


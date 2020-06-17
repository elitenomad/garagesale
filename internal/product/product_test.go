package product_test

import (
	"github.com/elitenomad/garagesale/internal/platform/database/databasetest"
	"github.com/elitenomad/garagesale/internal/product"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	np := product.NewProduct{
		Name:     "Comic Books",
		Cost:     10,
		Quantity: 20,
	}

	now := time.Date(2020, time.June, 1, 0, 0, 0, 0, time.UTC)
	p, err := product.Create(db, np, now)
	if err != nil {
		t.Fatalf("Could not create product: %v ", err)
	}

	saved, err := product.Fetch(db, p.ID)
	if err != nil {
		t.Fatalf("Could not retrieve the product %v ", err)
	}

	if diff := cmp.Diff(p, saved); diff != "" {
		t.Fatalf("Saved product did not match created: see diff %v", diff)
	}
}

func TestList(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	np_1 := product.NewProduct{
		Name:     "Comic Books",
		Cost:     10,
		Quantity: 20,
	}

	now := time.Date(2020, time.June, 1, 0, 0, 0, 0, time.UTC)
	// Alos we can seed database
	// schema.Seed(db)
	_, err := product.Create(db, np_1, now)
	if err != nil {
		t.Fatalf("Could not create product: %v ", err)
	}

	saved, err := product.List(db)
	if err != nil {
		t.Fatalf("Could not retrieve the product %v ", err)
	}

	if exp, got := 1, len(saved); exp != got {
		t.Fatalf("Saved products did not match created: see diff %v", got)
	}
}

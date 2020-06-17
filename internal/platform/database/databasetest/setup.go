package databasetest

import (
	"github.com/elitenomad/garagesale/internal/platform/database"
	"github.com/elitenomad/garagesale/internal/schema"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

func Setup(t *testing.T) (*sqlx.DB, func())  {
	t.Helper()

	// Start a docker container
	c := StartContainer(t)

	// Because i used docker to create postgres
	// These are dummy UN's and Passwords
	// In test ENV fetch this from environment variables.
	db, err := database.OpenDB(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil{
		t.Fatalf("Opening database connection %v", err)
	}

	t.Log("Waiting for database to be Ready...")

	var pingError error
	maxAttempts := 20

	for attempts := 1; attempts < maxAttempts; attempts++ {
		pingError  = db.Ping()
		if pingError == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		StopContainer(t, c)
		t.Fatalf("Wating for datbase to be ready : %v ", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		StopContainer(t, c)
		t.Fatalf("Migrating ... %v", err)
	}

	// Run the function passed
	teardown := func() {
		t.Helper()
		db.Close()
		StopContainer(t, c)
	}

	return db, teardown
}

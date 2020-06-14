package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
)

// Config is the configuration require to start the database
type Config struct {
	User string
	Password string
	Host string
	Name string
	DisableTLS bool
}

func OpenDB(cfg Config) (*sqlx.DB, error) {

	// Customise SSLMode
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := url.Values{}
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}


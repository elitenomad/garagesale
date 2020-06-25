package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/elitenomad/garagesale/internal/platform/conf"
	"github.com/elitenomad/garagesale/internal/platform/database"
	"github.com/elitenomad/garagesale/internal/schema"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:test,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	/*
		--------------------------------------------
		Open DB connection
		--------------------------------------------
	*/
	db, err := database.OpenDB(database.Config{
		User:       cfg.DB.User,
		Password:   "",
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "Database connection error")
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "Applying Migrations...")
		}
		log.Println("Migration complete...")
		return nil
	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "Applying Seed...")
		}
		log.Println("Seeding complete...")
		return nil
	}

	return nil
}

package main

import (
	"flag"
	"fmt"
	"github.com/elitenomad/garagesale/internal/platform/conf"
	"github.com/elitenomad/garagesale/internal/platform/database"
	"github.com/elitenomad/garagesale/internal/schema"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {

	// =========================================================================
	// Configuration

	var cfg struct {
		DB struct {
			User       string `conf:"default:pranavaswaroop"`
			Password   string `conf:"default:test,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	/*
		--------------------------------------------
		Open DB connection
		--------------------------------------------
	*/
	db, err := database.OpenDB(database.Config{
		User:      cfg.DB.User,
		Password:  "",
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("Applying Migrations...", err)
			os.Exit(1)
		}
		log.Println("Migration complete...")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("Applying Seed...", err)
			os.Exit(1)
		}
		log.Println("Seeding complete...")
		return
	}
}

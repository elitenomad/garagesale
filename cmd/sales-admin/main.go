package main

import (
	"flag"
	"github.com/elitenomad/garagesale/internal/schema"
	"github.com/elitenomad/garagesale/internal/platform/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	/*
		--------------------------------------------
		Open DB connection
		--------------------------------------------
	*/
	db, err := database.OpenDB()
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

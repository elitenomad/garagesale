package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/elitenomad/garagesale/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "fmt"

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("pranavaswaroop", ""),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}


func main() {
	/*
		--------------------------------------------
		App starting
		--------------------------------------------
	 */
	log.Println("Main started ...")
	defer log.Println("Main completed ...")

	/*
		--------------------------------------------
		Open DB connection
		--------------------------------------------
	*/
	db, err := openDB()
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

	service := Products{db: db}

	/*
		--------------------------------------------
		Starting API Service
		--------------------------------------------
	*/
	api := http.Server{
		Addr:              "localhost:3000",
		Handler:           http.HandlerFunc(service.List),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	/*
		--------------------------------------------
		Make a channel listen for errors coming from
		listener. Use a buffer channel so the goroutine
		can exit if we don't collect the error.
		--------------------------------------------
	*/
	serverErrors := make(chan error, 1)

	/*
		--------------------------------------------
		Start the service listening for requests.
		--------------------------------------------
	*/
	go func() {
		log.Printf("API is listening on %s ", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	/*
		--------------------------------------------
		Make a channel to listen for interruptions.
		--------------------------------------------
	*/
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)


	/*
		--------------------------------------------
		Wait before shutdown.
		--------------------------------------------
	*/
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)
	case <-interruption:
		fmt.Println("Main: starting shutdown")

		/*
			--------------------------------------------
			Set time out for completion.
			--------------------------------------------
		*/
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		/*
			--------------------------------------------
			Ask listener to shutdown.
			--------------------------------------------
		*/
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type Products struct {
	db *sqlx.DB
}

func (p *Products) List(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	const q = `SELECT * from products`

	if err := p.db.Select(&products, q); err != nil {
		log.Printf("error: selecting products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(products)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}


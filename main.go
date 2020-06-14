package main


/*
	TODO: This file is not used any more as we moved the functionality to packages.
		Remove the file once the project is over.
 */

//
//import (
//	"context"
//	"encoding/json"
//	"flag"
//	"github.com/elitenomad/garagesale/internal/schema"
//	"github.com/elitenomad/garagesale/internal/platform/database"
//	"github.com/jmoiron/sqlx"
//	_ "github.com/lib/pq"
//	"log"
//	"net/http"
//	"os"
//	"os/signal"
//	"syscall"
//	"time"
//)
//import "fmt"
//
//
//func main() {
//	/*
//		--------------------------------------------
//		App starting
//		--------------------------------------------
//	 */
//	log.Println("Main started ...")
//	defer log.Println("Main completed ...")
//
//	/*
//		--------------------------------------------
//		Open DB connection
//		--------------------------------------------
//	*/
//	db, err := database.OpenDB()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	flag.Parse()
//
//	switch flag.Arg(0) {
//	case "migrate":
//		if err := schema.Migrate(db); err != nil {
//			log.Println("Applying Migrations...", err)
//			os.Exit(1)
//		}
//		log.Println("Migration complete...")
//		return
//	case "seed":
//		if err := schema.Seed(db); err != nil {
//			log.Println("Applying Seed...", err)
//			os.Exit(1)
//		}
//		log.Println("Seeding complete...")
//		return
//	}
//
//	service := Products{db: db}
//
//	/*
//		--------------------------------------------
//		Starting API Service
//		--------------------------------------------
//	*/
//	api := http.Server{
//		Addr:              "localhost:3000",
//		Handler:           http.HandlerFunc(service.List),
//		ReadTimeout:       5 * time.Second,
//		WriteTimeout:      5 * time.Second,
//	}
//
//	/*
//		--------------------------------------------
//		Make a channel listen for errors coming from
//		listener. Use a buffer channel so the goroutine
//		can exit if we don't collect the error.
//		--------------------------------------------
//	*/
//	serverErrors := make(chan error, 1)
//
//	/*
//		--------------------------------------------
//		Start the service listening for requests.
//		--------------------------------------------
//	*/
//	go func() {
//		log.Printf("API is listening on %s ", api.Addr)
//		serverErrors <- api.ListenAndServe()
//	}()
//
//	/*
//		--------------------------------------------
//		Make a channel to listen for interruptions.
//		--------------------------------------------
//	*/
//	interruption := make(chan os.Signal, 1)
//	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)
//
//
//	/*
//		--------------------------------------------
//		Wait before shutdown.
//		--------------------------------------------
//	*/
//	select {
//	case err := <-serverErrors:
//		log.Fatalf("error: listening and serving: %s", err)
//	case <-interruption:
//		fmt.Println("Main: starting shutdown")
//
//		/*
//			--------------------------------------------
//			Set time out for completion.
//			--------------------------------------------
//		*/
//		const timeout = 5 * time.Second
//		ctx, cancel := context.WithTimeout(context.Background(), timeout)
//		defer cancel()
//
//		/*
//			--------------------------------------------
//			Ask listener to shutdown.
//			--------------------------------------------
//		*/
//		err := api.Shutdown(ctx)
//		if err != nil {
//			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
//			err = api.Close()
//		}
//
//		if err != nil {
//			log.Fatalf("main : could not stop server gracefully : %v", err)
//		}
//	}
//}
//
//
//type Products struct {
//	db *sqlx.DB
//}
//
//

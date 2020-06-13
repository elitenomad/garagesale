package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "fmt"

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
		Starting API Service
		--------------------------------------------
	*/
	api := http.Server{
		Addr:              "localhost:3000",
		Handler:           http.HandlerFunc(ListProducts),
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
	Name	string `json:"name"`
	Cost 	float64 `json:"cost"`
	Quantity int64 `json:"quantity"`
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	// This is not the real world scenario where we hardcode values
	// products := []Product{}
	// enchoding json will treat that as empty array rather than nil
	// This helps FE to handle bugs efficiently.
	products := []Product{
		{ Name: "Geography books", Cost: 200.25, Quantity: 10},
		{ Name: "History books", Cost: 300.25, Quantity: 15},
		{ Name: "Civics books", Cost: 289.25, Quantity: 20},
		{ Name: "Economics books", Cost: 100.25, Quantity: 25},
		{ Name: "Physics books", Cost: 220.25, Quantity: 30},
	}

	data, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// We can explicitly set the header status Ok but it should always be set
	// After the Header().set statements or else WriteHeader will override the values
	// w.WriteHeader(http.StausOK)
	if _, err := w.Write(data); err != nil {
		log.Println("Error in writing the data ", err)
	}
}


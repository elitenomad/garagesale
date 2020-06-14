package main

import (
	"context"
	"github.com/elitenomad/garagesale/cmd/sales-api/internal/handlers"
	"github.com/elitenomad/garagesale/internal/platform/conf"
	_ "github.com/lib/pq"
	"github.com/elitenomad/garagesale/internal/platform/database"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "fmt"


func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {
	// Global variable for this program.
	log := log.New(os.Stdout, "SALES: ", log.LstdFlags|log.Lmicroseconds|log.Llongfile)

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
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
				errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		errors.Wrap(err,"parsing config")
	}

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
	db, err := database.OpenDB(database.Config{
		User:      cfg.DB.User,
		Password:  "",
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		errors.Wrap(err,"Database connection issue")
	}
	defer db.Close()

	/*
		--------------------------------------------
		Starting API Service
		--------------------------------------------
	*/
	api := http.Server{
		Addr:              "localhost:8000",
		Handler:           handlers.API(log, db),
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
		errors.Wrap(err,"listening and serving error...")
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
			errors.Wrap(err, "main : could not stop server gracefully")
		}
	}

	return nil
}

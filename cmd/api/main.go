package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type app struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 5000, "Your server port")
	flag.StringVar(&cfg.env, "env", "dev", "Your environment: dev|prod|stage")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &app{
		config: cfg,
		logger: logger,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	err := server.ListenAndServe()
	logger.Fatal(err)
}

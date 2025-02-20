package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var conf config

	// cmd-line options
	flag.IntVar(&conf.port, "port", 8000, "Application port.")
	flag.StringVar(&conf.env, "env", "dev", "Application environment.")
	flag.Parse()

	// set the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: conf,
		logger: logger,
	}

	// define server settings
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server...", "addr", server.Addr, "env", conf.env)

	// start server
	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

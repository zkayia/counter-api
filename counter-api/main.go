package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var dbPath = flag.String(
	"db",
	envWithDefaultString("COUNTERAPI_DB_PATH", "database/"),
	"Path of the database, the process must have read/write permissions.",
)
var backupsPath = flag.String(
	"backups",
	envWithDefaultString("COUNTERAPI_BACKUPS_PATH", "backups/"),
	"Directory for database backups, the process must have read/write permissions.",
)
var logsPath = flag.String(
	"logs",
	envWithDefaultString("COUNTERAPI_LOGS_PATH", "logs/"),
	"Directory for log files, the process must have read/write permissions.",
)
var serverAddress = flag.String(
	"address",
	envWithDefaultString("COUNTERAPI_LISTEN_ADDRESS", "127.0.0.1:8000"),
	"Server address to listen on.",
)
var timeout = flag.Duration(
	"timeout",
	envWithDefaultDuration("COUNTERAPI_TIMEOUT", 15*time.Second),
	"Server timeout.",
)

func main() {

	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if setupLogFile() != nil {
		log.Println("[WARN] Failed to create log file, logging to stdout.")
	}

	log.Printf(`[INFO] Starting counter-api with:
	db=%s
	backups=%s
	logs=%s
	address=%s
	timeout=%s
	`,
		*dbPath,
		*backupsPath,
		*logsPath,
		*serverAddress,
		(*timeout).String(),
	)

	if err := dbCreateBackup(); err != nil {
		log.Printf("[WARN] Backup failed, %s\n", err)
	}
	ticker := time.NewTicker(48 * time.Hour)
	go func() {
		for range ticker.C {
			if err := dbCreateBackup(); err != nil {
				log.Printf("[WARN] Backup failed, %s\n", err)
			}
		}
	}()

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.URLFormat)
	mux.Use(middleware.Timeout(*timeout))
	mux.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))
	mux.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))

	counterArg := "{counter:[a-zA-Z0-9-_]+}"

	router := mux.With(counterContext)

	router.With(getQueryContext).Get("/get/"+counterArg, handleGet)

	for _, operation := range Operations {
		router.With(valueQueryContext).With(operationContext(operation)).Get(
			"/"+operation.toString()+"/"+counterArg,
			handleOperation,
		)
	}

	srv := &http.Server{
		Handler:      mux,
		Addr:         *serverAddress,
		WriteTimeout: *timeout,
		ReadTimeout:  *timeout,
	}

	log.Fatal(srv.ListenAndServe())
}

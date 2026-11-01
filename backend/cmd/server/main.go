package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"securecoda/internal/api"
	"securecoda/internal/coda"
	"securecoda/internal/detector"
	"securecoda/internal/scheduler"
	"securecoda/internal/store"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = time.RFC3339
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		logger.Warn().Err(err).Msg("No .env file found")
	}

	// Load env
	codaToken := os.Getenv("CODA_API_TOKEN")
	if codaToken == "" {
		logger.Fatal().Msg("CODA_API_TOKEN is required")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "securecoda.db"
	}

	// Initialize store
	st, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init store")
	}

	// Initialize Coda client
	codaClient := coda.NewClient(codaToken, logger)

	// Detector
	det := detector.NewDetector(codaClient, st, logger)

	// Scheduler (poller)
	pollInterval := 60 * time.Second
	if v := os.Getenv("POLL_INTERVAL_SEC"); v != "" {
		// ignore parse errors for brevity
	}
	sched := scheduler.NewScheduler(det, pollInterval, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sched.Start(ctx)

	// API Server
	r := mux.NewRouter()
	api.RegisterRoutes(r, st, det, codaClient, logger)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	logger.Info().Msg("SecureCoda backend starting on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("server error")
	}
}

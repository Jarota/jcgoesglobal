package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/jarota/jctravels/internal/handler"
	"github.com/jarota/jctravels/internal/storage"
)

var (
	dbDSN     string
	uploadDir string
	httpPort  string
)

func main() {
	err := loadEnv()
	if err != nil {
		slog.Error("failed to load env", slog.Any("err", err))
		return
	}

	slog.Info("initialising storage...")
	store, err := storage.New(dbDSN, uploadDir)
	if err != nil {
		slog.Error("failed to create new store", slog.Any("err", err))
		return
	}
	defer func() {
		err = store.Close()
		if err != nil {
			slog.Error("failed to close store", slog.Any("err", err))
		}
	}()

	slog.Info("registering endpoint handlers...")
	h := handler.New(store)
	http.HandleFunc("POST /api/new", h.NewPost)
	http.HandleFunc("GET /api/all", h.AllPosts)

	slog.Info("listening for requests...")
	err = http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		slog.Error("listen and serve errored", slog.Any("err", err))
	}
}

func loadEnv() error {
	httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return errors.New("missing http port env var")
	}

	dbDSN = os.Getenv("DB_DSN")
	if dbDSN == "" {
		return errors.New("missing db dsn env var")
	}

	uploadDir = os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		return errors.New("missing upload dir env var")
	}

	return nil
}

package main

import (
	"log/slog"
	"net/http"

	"github.com/jarota/jctravels/internal/handler"
	"github.com/jarota/jctravels/internal/storage"
)

func main() {

	slog.Info("initialising storage...")
	store, err := storage.New("file:app.db")
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

	slog.Info("listing for requests...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("listen and serve errored", slog.Any("err", err))
	}
}

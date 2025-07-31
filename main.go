package main

import (
	"log/slog"
	"net/http"

	"github.com/jarota/jctravels/internal/handler"
	"github.com/jarota/jctravels/internal/storage"
)

func main() {

	store, err := storage.New("localhost")
	if err != nil {
		slog.Error("failed to create new store", slog.Any("err", err))
		return
	}

	h := handler.New(store)
	http.HandleFunc("POST /api/new", h.NewPost)
	http.HandleFunc("GET /api/all", h.AllPosts)

	slog.Info("handlers registered, starting server...")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("listen and serve errored", slog.Any("err", err))
	}
}

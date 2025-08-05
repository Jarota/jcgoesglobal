package handler

import (
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/jarota/jctravels/internal/model"
)

const (
	FormCaptionKey = "caption"
	FormFileKey    = "image-upload"
)

type Store interface {
	CreatePost(postID, caption, author string, likes int) error
	CreateImages(postID string, fileHeaders []*multipart.FileHeader) error
	LookupAll() ([]*model.Post, error)
}

type handler struct {
	store Store
}

func New(store Store) *handler {
	return &handler{store}
}

func (h *handler) NewPost(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling request to new post")

	author, _, ok := r.BasicAuth()
	if !ok {
		slog.Error("missing basic auth details")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		slog.Error("failed to parse form", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postID := uuid.NewString()
	caption := r.MultipartForm.Value[FormCaptionKey][0]
	err = h.store.CreatePost(postID, caption, author, 0)
	if err != nil {
		slog.Error("failed to create post", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Save any images uploaded too
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		err = h.store.CreateImages(postID, r.MultipartForm.File[FormFileKey])
		if err != nil {
			slog.Error("failed to create images", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *handler) AllPosts(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling request to all posts")

	posts, err := h.store.LookupAll()
	if err != nil {
		slog.Error("failed to lookup all posts", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		slog.Error("failed to encode posts", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

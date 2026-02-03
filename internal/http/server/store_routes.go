package server

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StoreRoutes(db *db.MemDB) http.Handler {
	r := chi.NewRouter()

	r.Post("/", handler.Store(db))

	return r
}

package server

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UpdateRoutes(db db.DB) http.Handler {
	r := chi.NewRouter()

	r.Put("/", handler.Update(db))

	return r
}

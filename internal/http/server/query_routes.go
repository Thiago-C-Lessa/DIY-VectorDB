package server

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func FetchRoutes(db db.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/all", handler.Query_All(db))
	r.Get("/one/", handler.Select(db))
	r.Get("/similar/", handler.SelectSimilar(db))

	return r
}

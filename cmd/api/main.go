package main

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/http/server"

	"log"
	"net/http"
)

func main() {
	dbm := db.NewMemDB()

	r := server.NewRouter()

	r.Mount("/store", server.StoreRoutes(dbm))
	r.Mount("/fetch", server.FetchRoutes(dbm))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

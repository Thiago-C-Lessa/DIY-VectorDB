package main

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/http/server"

	"log"
	"net/http"
)

func main() {
	dbvm := db.NewVecMemDB()

	r := server.NewRouter()

	r.Mount("/store", server.StoreRoutes(dbvm))
	r.Mount("/fetch", server.FetchRoutes(dbvm))
	r.Mount("/update", server.UpdateRoutes(dbvm))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package handler

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/models"
	"encoding/json"
	"net/http"
)

func Query_All(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.ListAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Select(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonKey models.RequestFetchOne
		err := json.NewDecoder(r.Body).Decode(&jsonKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := db.Select(jsonKey.Key)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func SelectSimilar(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonKey models.RequestFetchSimilar
		err := json.NewDecoder(r.Body).Decode(&jsonKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := db.SelectSimilar(jsonKey.Key, jsonKey.K)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

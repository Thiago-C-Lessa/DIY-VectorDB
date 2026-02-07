package handler

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/models"
	"encoding/json"
	"net/http"
)

func Store(db db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rawData map[string]json.RawMessage
		var rData models.RequestPostData

		if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, ok := rawData["key"]
		if !ok {
			http.Error(w, "campo 'key' é obrigatório", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(rawData["key"], &rData.Key); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		delete(rawData, "key") //apaga a key e só deixa o conteúdo da requisição
		var err error
		rData.Rest, err = json.Marshal(rawData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.Insert(rData.Key, rData.Rest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(rData.Rest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

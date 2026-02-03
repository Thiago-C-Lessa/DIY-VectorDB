package db

import (
	"DIY-VectorDB/internal/exceptions"
	"DIY-VectorDB/internal/models"
	"encoding/json"
	"sync"

	"DIY-VectorDB/internal/http/embedding"
)

type MemDB struct {
	sync.RWMutex
	embeddings map[string][768]float32
	data       map[[768]float32]json.RawMessage
}

func NewMemDB() *MemDB {
	return &MemDB{
		embeddings: make(map[string][768]float32),
		data:       make(map[[768]float32]json.RawMessage),
	}
}

func (db *MemDB) Insert(key string, value json.RawMessage) error {
	db.Lock()
	defer db.Unlock()

	_, ok := db.embeddings[key]
	if ok {
		return &exceptions.ErrorKeyAlreadyExists{Key: key}
	}

	emb, err := embedding.GenereteEmbedding(key)
	if err != nil {
		return err
	}

	db.embeddings[key] = emb
	db.data[emb] = value
	return nil

}

func (db *MemDB) ListAll() (models.ResponseData, error) {
	db.RLock()
	defer db.RUnlock()

	var data models.ResponseData

	for key, emb := range db.embeddings {
		content, ok := db.data[emb]
		if !ok {
			return data, &exceptions.ErrorContentNotFound{Key: key}
		}

		data.Keys = append(data.Keys, key)
		data.Contents = append(data.Contents, content)
	}

	return data, nil
}

func (db *MemDB) Select(key string) (models.ResponseData, error) {
	db.RLock()
	defer db.RUnlock()

	var data models.ResponseData

	emb, ok := db.embeddings[key]
	if !ok {
		return data, &exceptions.ErrorContentNotFound{Key: key}
	}

	d, _ := db.data[emb]
	data.Keys = append(data.Keys, key)
	data.Contents = append(data.Contents, d)

	return data, nil
}

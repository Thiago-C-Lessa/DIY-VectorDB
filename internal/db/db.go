package db

import (
	"DIY-VectorDB/internal/models"
	"encoding/json"
)

type DB interface {
	Insert(key string, value json.RawMessage) error
	ListAll() (models.ResponseData, error)
	Select(key string) (models.ResponseData, error)
	//SelectSimilar(key [768]float64) ([]interface{}, error)
	//Delete(key [768]float64) error
}

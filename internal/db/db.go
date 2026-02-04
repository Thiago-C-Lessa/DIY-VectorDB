package db

import (
	"DIY-VectorDB/internal/models"
	"encoding/json"
)

type DB interface {
	Insert(key string, value json.RawMessage) error
	ListAll() (models.ResponseData, error)
	Select(key string) (models.ResponseData, error)
	SelectSimilar(key string, k uint64) (models.ResponseData, error)
	//Delete(key string) error
	Update(key string, value json.RawMessage) error
}

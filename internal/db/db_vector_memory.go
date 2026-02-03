package db

import (
	"encoding/json"
	"log/slog"
	"sync"
)

type node struct {
	embedding       [768]float32
	level uint64
	neighbors [][]*node
}

type VecMemDB struct {
	sync.RWMutex
	maxLevel uint64
	embeddings map[string][768]float32
	data       map[[768]float32]json.RawMessage
	levels [][]*node
}

func NewVecMemDB() *VecMemDB {
	return &VecMemDB{
		embeddings: make(map[string][768]float32),
		data:       make(map[[768]float32]json.RawMessage),
	}
}

func newNode(key [768]float32) *node {
	n := new(node)

	n.embedding = key
	return n
}


func
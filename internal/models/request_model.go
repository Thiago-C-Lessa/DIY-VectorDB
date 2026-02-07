package models

import "encoding/json"

type RequestPostData struct {
	Key  string          `json:"key"`
	Rest json.RawMessage `json:"-"`
}

type RequestFetchOne struct {
	Key string `json:"key"`
}

type RequestFetchSimilar struct {
	Key string `json:"key"`
	K   uint64 `json:"k"`
}

type RequestPut struct {
	Key  string          `json:"key"`
	Rest json.RawMessage `json:"-"`
}

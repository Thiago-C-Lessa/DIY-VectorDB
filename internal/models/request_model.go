package models

import "encoding/json"

type RequestData struct {
	Key  string          `json:"key"`
	Rest json.RawMessage `json:"-"`
}

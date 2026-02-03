package models

import "encoding/json"

type ResponseData struct {
	Keys     []string          `json:"keys"`
	Contents []json.RawMessage `json:"contents"`
}

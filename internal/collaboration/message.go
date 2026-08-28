package collaboration

import "encoding/json"

type Message struct {
	Type         string          `json:"type"`
	ProjectID    string          `json:"projectId"`
	UserID       string          `json:"userId"`
	Payload      json.RawMessage `json:"payload"`
	BaseRevision int             `json:"baseRevision"`
}
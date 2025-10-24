package websocket

import "encoding/json"

type MessageWs struct {
	Type    string          `json:"type"` // e.g. "chat", "notify", "ping"
	From    string          `json:"from,omitempty"`
	To      string          `json:"to,omitempty"` // optional target userID
	Payload json.RawMessage `json:"payload,omitempty"`
}

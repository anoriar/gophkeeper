package response

import (
	"encoding/json"
)

type CommandResponse struct {
	Status  string          `json:"status"`
	Error   string          `json:"error"`
	Payload json.RawMessage `json:"payload"`
}

package command

type CommandResponse struct {
	Status  string      `json:"status"`
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
}

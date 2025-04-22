package response

type Error struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

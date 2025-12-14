package reply

type ReplyUpdateDTO struct {
	Body     string `json:"body,omitempty"`
	Quantity uint64 `json:"quantity,omitempty"`
}

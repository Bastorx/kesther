package model

// Callback represents a JSON input/output payload of an event callback
type Callback struct {
	ID     string `json:"id"`
	PlanID string `json:"planId"`
	Title  string `json:"title"`
	Undo   Action `json:"undo"`
	Parent Action `json:"parent"`
}

// Action represents a possible action (apply to parent, or undo)
type Action struct {
	PlanID  string      `json:"planId"`
	Method  string      `json:"method"`
	URI     string      `json:"uri"`
	Payload interface{} `json:"payload"`
}

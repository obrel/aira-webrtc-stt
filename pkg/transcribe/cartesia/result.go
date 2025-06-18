package cartesia

type Result struct {
	Type    string `json:"type"`
	IsFinal bool   `json:"is_final"`
	Text    string `json:"text"`
}

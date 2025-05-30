package deepgram

type Result struct {
	Type    string  `json:"type"`
	IsFinal bool    `json:"is_final"`
	Channel Channel `json:"channel"`
}

type Channel struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript string  `json:"transcript"`
	Confidence float32 `json:"confidence"`
}

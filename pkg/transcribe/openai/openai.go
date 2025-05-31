package openai

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
)

type OpenAITranscriber struct {
	client *websocket.Conn
	ctx    context.Context
}

func NewOpenAI(ctx context.Context, apiKey string) (transcribe.Service, error) {
	websocket.DefaultDialer.Subprotocols = []string{
		"realtime",
		"openai-beta.realtime-v1",
		"openai-insecure-api-key." + apiKey,
	}

	client, _, err := websocket.DefaultDialer.Dial("wss://api.openai.com/v1/realtime?intent=transcription", nil)
	if err != nil {
		return nil, err
	}

	return &OpenAITranscriber{
		client: client,
		ctx:    ctx,
	}, nil
}

func (t *OpenAITranscriber) CreateStream() (transcribe.Stream, error) {
	return &OpenAIStream{
		stream:  t.client,
		results: make(chan transcribe.Result),
		mu:      sync.Mutex{},
	}, nil
}

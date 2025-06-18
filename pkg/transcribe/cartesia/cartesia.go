package cartesia

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
)

type CartesiaTranscriber struct {
	client *websocket.Conn
	ctx    context.Context
}

func NewCartesia(ctx context.Context, apiKey string) (transcribe.Service, error) {
	client, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://api.cartesia.ai/stt/websocket?cartesia_version=2024-11-13&model=ink-whisper&language=en&sample_rate=16000&encoding=pcm_s16le&api_key=%s", apiKey), nil)
	if err != nil {
		return nil, err
	}

	return &CartesiaTranscriber{
		client: client,
		ctx:    ctx,
	}, nil
}

func (t *CartesiaTranscriber) CreateStream() (transcribe.Stream, error) {
	return &CartesiaStream{
		stream:  t.client,
		results: make(chan transcribe.Result),
		mu:      sync.Mutex{},
	}, nil
}

package deepgram

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
)

type DeepgramTranscriber struct {
	client *websocket.Conn
	ctx    context.Context
}

func NewDeepgram(ctx context.Context, apiKey string) (transcribe.Service, error) {
	websocket.DefaultDialer.Subprotocols = []string{"token", apiKey}

	client, _, err := websocket.DefaultDialer.Dial("wss://api.deepgram.com/v1/listen?model=nova-2&language=id&encoding=linear16&sample_rate=48000", nil)
	if err != nil {
		return nil, err
	}

	return &DeepgramTranscriber{
		client: client,
		ctx:    ctx,
	}, nil
}

func (t *DeepgramTranscriber) CreateStream() (transcribe.Stream, error) {
	return &DeepgramStream{
		stream:  t.client,
		results: make(chan transcribe.Result),
	}, nil
}

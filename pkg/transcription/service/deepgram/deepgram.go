package deepgram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"github.com/obrel/go-lib/pkg/log"
)

const baseURL = "wss://api.deepgram.com/v1/listen"

type Deepgram struct {
	apiKey     string
	dialURL    string
	model      string
	language   string
	encoding   string
	sampleRate int
	wsClient   *websocket.Conn
	lock       sync.Mutex
}

func (d *Deepgram) Connect() error {
	var err error
	var res *http.Response

	websocket.DefaultDialer.Subprotocols = []string{"token", d.apiKey}

	d.wsClient, res, err = websocket.DefaultDialer.Dial(d.dialURL, nil)
	if err != nil {
		log.For("deepgram", "connect").Error(res.Status)
		return err
	}

	return nil
}

func (d *Deepgram) Write(stream []byte) (int, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if err := d.wsClient.WriteMessage(websocket.BinaryMessage, stream); err != nil {
		return 0, err
	}

	return len(stream), nil
}

func (d *Deepgram) Receive(res chan transcription.Result, done chan bool) error {
	for {
		select {
		case <-done:
			return nil
		default:
			_, resp, err := d.wsClient.ReadMessage()
			if err != nil && err != io.EOF {
				return err
			}

			result := &Result{}

			err = json.Unmarshal(resp, result)
			if err != nil {
				log.For("deepgram", "receive").Error(err)
				continue
			}

			for _, alt := range result.Channel.Alternatives {
				if alt.Transcript != "" {
					log.Printf("%s (%.2f)", alt.Transcript, alt.Confidence)

					res <- transcription.Result{
						Confidence: alt.Confidence,
						Text:       alt.Transcript,
						Final:      result.IsFinal,
					}
				}
			}
		}
	}
}

func (d *Deepgram) Close() error {
	_ = d.wsClient.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"Close\" }"))
	return d.wsClient.Close()
}

func init() {
	transcription.Register("deepgram", func(opts ...transcription.Option) (transcription.Transcription, error) {
		s := &Deepgram{
			lock: sync.Mutex{},
		}

		for _, opt := range opts {
			switch f := opt.(type) {
			case func(*Deepgram):
				f(s)
			default:
				return nil, fmt.Errorf("Unknown option.")
			}
		}

		s.dialURL = fmt.Sprintf("%s?model=%s&language=%s&encoding=%s&sample_rate=%v",
			baseURL,
			s.model,
			s.language,
			s.encoding,
			s.sampleRate,
		)

		return s, nil
	})
}

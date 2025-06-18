package cartesia

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"github.com/obrel/go-lib/pkg/log"
)

const baseURL = "wss://api.cartesia.ai/stt/websocket?cartesia_version=2024-11-13"

type Cartesia struct {
	apiKey     string
	model      string
	language   string
	encoding   string
	sampleRate int
	wsClient   *websocket.Conn
	lock       sync.Mutex
}

func (c *Cartesia) Write(stream []byte) (int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := c.wsClient.WriteMessage(websocket.BinaryMessage, stream); err != nil {
		return 0, err
	}

	return len(stream), nil
}

func (c *Cartesia) Receive(res chan transcription.Result, done chan bool) error {
	for {
		select {
		case stop := <-done:
			fmt.Println("BAJING")
			if stop {
				return nil
			}
		default:
			_, resp, err := c.wsClient.ReadMessage()
			if err != nil && err != io.EOF {
				return err
			}

			result := &Result{}

			err = json.Unmarshal(resp, result)
			if err != nil {
				log.For("cartesia", "receive").Error(err)
				continue
			}

			if result.Type == "transcript" && result.Text != "" {
				log.Printf(result.Text)

				res <- transcription.Result{
					Text:  result.Text,
					Final: result.IsFinal,
				}
			}
		}
	}
}

func (c *Cartesia) Close() error {
	return c.wsClient.Close()
}

func init() {
	transcription.Register("cartesia", func(opts ...transcription.Option) (transcription.Transcription, error) {
		var err error
		s := &Cartesia{
			lock: sync.Mutex{},
		}

		for _, opt := range opts {
			switch f := opt.(type) {
			case func(*Cartesia):
				f(s)
			default:
				return nil, fmt.Errorf("Unknown option.")
			}
		}

		dialURL := fmt.Sprintf("%s&api_key=%s&model=%s&language=%s&encoding=%s&sample_rate=%v",
			baseURL,
			s.apiKey,
			s.model,
			s.language,
			s.encoding,
			s.sampleRate,
		)

		s.wsClient, _, err = websocket.DefaultDialer.Dial(dialURL, nil)
		if err != nil {
			return nil, err
		}

		return s, nil
	})
}

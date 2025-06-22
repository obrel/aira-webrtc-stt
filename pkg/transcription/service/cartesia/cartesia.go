package cartesia

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

const baseURL = "wss://api.cartesia.ai/stt/websocket?cartesia_version=2024-11-13"

type Cartesia struct {
	apiKey     string
	dialURL    string
	model      string
	language   string
	encoding   string
	sampleRate int
	wsClient   *websocket.Conn
	lock       sync.Mutex
}

func (c *Cartesia) Connect() error {
	var err error
	var res *http.Response

	c.wsClient, res, err = websocket.DefaultDialer.Dial(c.dialURL, nil)
	if err != nil {
		log.For("cartesia", "connect").Error(res.Status)
		return err
	}

	return nil
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
		case <-done:
			return nil
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
	_ = c.wsClient.WriteMessage(websocket.TextMessage, []byte("done"))
	return c.wsClient.Close()
}

func init() {
	transcription.Register("cartesia", func(opts ...transcription.Option) (transcription.Transcription, error) {
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

		s.dialURL = fmt.Sprintf("%s&api_key=%s&model=%s&language=%s&encoding=%s&sample_rate=%v",
			baseURL,
			s.apiKey,
			s.model,
			s.language,
			s.encoding,
			s.sampleRate,
		)

		return s, nil
	})
}

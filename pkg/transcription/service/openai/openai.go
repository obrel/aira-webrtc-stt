package openai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"github.com/obrel/go-lib/pkg/log"
)

const baseURL = "wss://api.openai.com/v1/realtime?intent=transcription"

type OpenAI struct {
	apiKey                   string
	dialURL                  string
	model                    string
	language                 string
	encoding                 string
	sampleRate               int
	prompt                   string
	detectionType            string
	detectionTheshold        float32
	detectionPrefixPadding   int
	detectionSilenceDuration int
	noiseReductionType       string
	wsClient                 *websocket.Conn
	ready                    bool
	lock                     sync.Mutex
}

type Data struct {
	Type  string `json:"type"`
	Audio string `json:"audio"`
}

func (o *OpenAI) Connect() error {
	var err error
	var res *http.Response

	websocket.DefaultDialer.Subprotocols = []string{
		"realtime",
		"openai-beta.realtime-v1",
		"openai-insecure-api-key." + o.apiKey,
	}

	o.wsClient, res, err = websocket.DefaultDialer.Dial(o.dialURL, nil)
	if err != nil {
		log.For("openai", "connect").Error(res.Status)
		return err
	}

	return nil
}

func (o *OpenAI) Write(buffer []byte) (int, error) {
	if !o.ready {
		return 0, nil
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	data := Data{
		Type:  "input_audio_buffer.append",
		Audio: base64.StdEncoding.EncodeToString(buffer),
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	if err := o.wsClient.WriteMessage(websocket.TextMessage, raw); err != nil {
		return 0, err
	}

	return len(buffer), nil
}

func (o *OpenAI) Send(buffer []byte) (int, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if err := o.wsClient.WriteMessage(websocket.TextMessage, buffer); err != nil {
		return 0, err
	}

	return len(buffer), nil
}

func (o *OpenAI) Receive(res chan transcription.Result, done chan bool) error {
	for {
		select {
		case <-done:
			return nil
		default:
			_, resp, err := o.wsClient.ReadMessage()
			if err != nil && err != io.EOF {
				return err
			}

			result := &Result{}

			err = json.Unmarshal(resp, result)
			if err != nil {
				log.For("openai", "receive").Error(err)
				continue
			}

			if result.Type == "transcription_session.created" {
				if !o.ready {
					msg := fmt.Sprintf(
						"{ \"type\": \"transcription_session.update\", \"session\": { \"input_audio_format\": \"%s\", \"input_audio_transcription\": { \"model\": \"%s\", \"language\": \"%s\", \"prompt\": \"%s\" }, \"turn_detection\": { \"type\": \"%s\", \"threshold\": %v, \"prefix_padding_ms\": %v, \"silence_duration_ms\": %v }, \"input_audio_noise_reduction\": { \"type\": \"%s\" } } }",
						o.encoding,
						o.model,
						o.language,
						o.prompt,
						o.detectionType,
						o.detectionTheshold,
						o.detectionPrefixPadding,
						o.detectionSilenceDuration,
						o.noiseReductionType,
					)

					err := o.wsClient.WriteMessage(websocket.TextMessage, []byte(msg))
					if err != nil {
						log.For("openai", "receive").Error(err)
					}

					o.ready = true
				}
			} else if result.Type == "conversation.item.input_audio_transcription.completed" {
				log.Printf(result.Transcript)

				// HACK: Sometimes it returns duplicate transcription with new line.
				trans := strings.Split(result.Transcript, "\n")

				res <- transcription.Result{
					Text:  trans[0],
					Final: true,
				}
			}
		}
	}
}

func (o *OpenAI) Close() error {
	o.ready = false
	_ = o.wsClient.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"input_audio_buffer.cleared\" }"))
	return o.wsClient.Close()
}

func init() {
	transcription.Register("openai", func(opts ...transcription.Option) (transcription.Transcription, error) {
		s := &OpenAI{
			ready: false,
			lock:  sync.Mutex{},
		}
		s.SetDefaultOptions()

		for _, opt := range opts {
			switch f := opt.(type) {
			case func(*OpenAI):
				f(s)
			default:
				return nil, fmt.Errorf("Unknown option.")
			}
		}

		s.dialURL = baseURL

		return s, nil
	})
}

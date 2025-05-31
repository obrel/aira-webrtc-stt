package openai

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
	"github.com/obrel/go-lib/pkg/log"
)

type Data struct {
	Type  string `json:"type"`
	Audio string `json:"audio"`
}

type OpenAIStream struct {
	stream  *websocket.Conn
	results chan transcribe.Result
	ready   bool
	mu      sync.Mutex
}

func (st *OpenAIStream) Write(buffer []byte) (int, error) {
	if !st.ready {
		return 0, nil
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	data := Data{
		Type:  "input_audio_buffer.append",
		Audio: base64.StdEncoding.EncodeToString(buffer),
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	if err := st.stream.WriteMessage(websocket.TextMessage, raw); err != nil {
		return 0, err
	}

	return len(buffer), nil
}

func (st *OpenAIStream) Recv(res chan transcribe.Result, done chan bool) error {
	for {
		select {
		case stop := <-done:
			if stop {
				return nil
			}
		default:
			_, resp, err := st.stream.ReadMessage()
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
				st.ready = true
				msg := `{"type":"transcription_session.update","session":{"input_audio_format":"pcm16","input_audio_transcription":{"model":"gpt-4o-transcribe","prompt":"","language":"id"},"turn_detection":{"type":"server_vad","threshold":0.5,"prefix_padding_ms":300,"silence_duration_ms":500},"input_audio_noise_reduction":{"type":"near_field"}}}`

				err := st.stream.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.For("openai", "receive").Error(err)
				}
			} else if result.Type == "conversation.item.input_audio_transcription.completed" {
				log.Printf(result.Transcript)

				// HACK: Sometimes it returns duplicate transcription with new line.
				trans := strings.Split(result.Transcript, "\n")

				res <- transcribe.Result{
					Text:  trans[0],
					Final: true,
				}
			}
		}
	}
}

func (st *OpenAIStream) Results() <-chan transcribe.Result {
	return st.results
}

func (st *OpenAIStream) Close() error {
	return st.stream.Close()
}

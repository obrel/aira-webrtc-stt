package deepgram

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
	"github.com/obrel/go-lib/pkg/log"
)

type DeepgramStream struct {
	stream  *websocket.Conn
	results chan transcribe.Result
	ready   bool
	mu      sync.Mutex
}

func (st *DeepgramStream) Write(buffer []byte) (int, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if err := st.stream.WriteMessage(websocket.BinaryMessage, buffer); err != nil {
		return 0, err
	}

	return len(buffer), nil
}

func (st *DeepgramStream) Recv(res chan transcribe.Result, done chan bool) error {
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
				log.For("deepgrap", "receive").Error(err)
				continue
			}

			for _, alt := range result.Channel.Alternatives {
				if alt.Transcript != "" {
					log.Printf("%s (%.2f)", alt.Transcript, alt.Confidence)

					res <- transcribe.Result{
						Confidence: alt.Confidence,
						Text:       alt.Transcript,
						Final:      result.IsFinal,
					}
				}
			}
		}
	}
}

func (st *DeepgramStream) Results() <-chan transcribe.Result {
	return st.results
}

func (st *DeepgramStream) Close() error {
	return st.stream.Close()
}

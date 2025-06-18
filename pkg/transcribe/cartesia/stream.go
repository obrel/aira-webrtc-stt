package cartesia

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
	"github.com/obrel/go-lib/pkg/log"
)

type CartesiaStream struct {
	stream  *websocket.Conn
	results chan transcribe.Result
	mu      sync.Mutex
}

func (st *CartesiaStream) Write(buffer []byte) (int, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if err := st.stream.WriteMessage(websocket.BinaryMessage, buffer); err != nil {
		return 0, err
	}

	return len(buffer), nil
}

func (st *CartesiaStream) Recv(res chan transcribe.Result, done chan bool) error {
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
				log.For("cartesia", "receive").Error(err)
				continue
			}

			if result.Type == "transcript" && result.Text != "" {
				log.Printf(result.Text)

				res <- transcribe.Result{
					Text:  result.Text,
					Final: result.IsFinal,
				}
			}
		}
	}
}

func (st *CartesiaStream) Results() <-chan transcribe.Result {
	return st.results
}

func (st *CartesiaStream) Close() error {
	return st.stream.Close()
}

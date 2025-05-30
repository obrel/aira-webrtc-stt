package google

import (
	"fmt"
	"io"
	"log"

	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type GoogleStream struct {
	stream  speechpb.Speech_StreamingRecognizeClient
	results chan transcribe.Result
}

func (st *GoogleStream) Write(buffer []byte) (int, error) {
	if err := st.stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
			AudioContent: buffer,
		},
	}); err != nil {
		return 0, nil
	}

	return len(buffer), nil
}

func (st *GoogleStream) Recv(res chan transcribe.Result, done chan bool) error {
	for {
		select {
		case stop := <-done:
			if stop {
				return nil
			}
		default:
			resp, err := st.stream.Recv()
			if err != nil && err != io.EOF {
				return err
			}

			if resp != nil {
				if resp.Error != nil {
					return fmt.Errorf("(Code: %d) %s", resp.Error.GetCode(), resp.Error.GetMessage())
				}

				for _, result := range resp.Results {
					for _, alt := range result.GetAlternatives() {
						log.Printf("%s (%.2f)", alt.GetTranscript(), alt.GetConfidence())

						res <- transcribe.Result{
							Confidence: alt.GetConfidence(),
							Text:       alt.GetTranscript(),
							Final:      result.GetIsFinal(),
						}
					}
				}
			}
		}
	}
}

func (st *GoogleStream) Results() <-chan transcribe.Result {
	return st.results
}

func (st *GoogleStream) Close() error {
	return st.stream.CloseSend()
}

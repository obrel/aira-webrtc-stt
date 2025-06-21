package google

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type Google struct {
	ctx          context.Context
	credential   string
	model        string
	language     string
	encoding     string
	sampleRate   int
	audioChannel int
	stream       speechpb.Speech_StreamingRecognizeClient
	lock         sync.Mutex
}

func (g *Google) Connect() error {
	client, err := speech.NewClient(g.ctx, option.WithCredentialsFile(g.credential))
	if err != nil {
		return err
	}

	g.stream, err = client.StreamingRecognize(g.ctx)
	if err != nil {
		return err
	}

	err = g.stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Model:             g.model,
					Encoding:          getEncoding(g.encoding),
					SampleRateHertz:   int32(g.sampleRate),
					LanguageCode:      g.language,
					AudioChannelCount: int32(g.audioChannel),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Google) Write(buffer []byte) (int, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	err := g.stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
			AudioContent: buffer,
		},
	})
	if err != nil {
		return 0, nil
	}

	return len(buffer), nil
}

func (g *Google) Receive(res chan transcription.Result, done chan bool) error {
	for {
		select {
		case <-done:
			return nil
		default:
			resp, err := g.stream.Recv()
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

						res <- transcription.Result{
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

func (g *Google) Close() error {
	return g.stream.CloseSend()
}

func init() {
	transcription.Register("google", func(opts ...transcription.Option) (transcription.Transcription, error) {
		s := &Google{
			ctx:  context.Background(),
			lock: sync.Mutex{},
		}

		for _, opt := range opts {
			switch f := opt.(type) {
			case func(*Google):
				f(s)
			default:
				return nil, fmt.Errorf("Unknown option.")
			}
		}

		return s, nil
	})
}

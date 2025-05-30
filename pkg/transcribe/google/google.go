package google

import (
	"context"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type GoogleTranscriber struct {
	speechClient *speech.Client
	ctx          context.Context
}

func (t *GoogleTranscriber) CreateStream() (transcribe.Stream, error) {
	stream, err := t.speechClient.StreamingRecognize(t.ctx)
	if err != nil {
		return nil, err
	}

	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:          speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz:   48000,
					LanguageCode:      "id",
					AudioChannelCount: 1,
				},
			},
		},
	}); err != nil {
		return nil, err
	}

	return &GoogleStream{
		stream:  stream,
		results: make(chan transcribe.Result),
	}, nil
}

func NewGoogleSpeech(ctx context.Context, credentials string) (transcribe.Service, error) {
	speechClient, err := speech.NewClient(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		return nil, err
	}
	return &GoogleTranscriber{
		speechClient: speechClient,
		ctx:          ctx,
	}, nil
}

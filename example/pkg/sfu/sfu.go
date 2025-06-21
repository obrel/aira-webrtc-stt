package sfu

import (
	"github.com/obrel/aira-websocket-stt/example/pkg/transcriber"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/pion/webrtc/v4"
)

type SFU struct {
	transcription transcription.Transcription
}

func NewSFU(transcription transcription.Transcription) *SFU {
	return &SFU{
		transcription: transcription,
	}
}

func (s *SFU) CreatePeerConnection() (PeerConnection, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, err
	}

	dataChan := make(chan *webrtc.DataChannel)
	stop := make(chan bool, 1)

	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})
	if err != nil {
		log.For("sfu", "peer").Error(err)
		return nil, err
	}

	pc.OnTrack(func(track *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		log.Printf("Received audio (%s) track, id = %s\n", track.Codec().MimeType, track.ID())

		err := transcriber.Transcribe(s.transcription, 16000, track, <-dataChan, stop)
		if err != nil {
			log.For("sfu", "peer").Fatal(err)
		}
	})

	pc.OnICEConnectionStateChange(func(connState webrtc.ICEConnectionState) {
		log.Printf("Connection state: %s \n", connState.String())

		if connState == webrtc.ICEConnectionStateClosed {
			stop <- true
		}
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dataChan <- dc
	})

	return &Client{
		pc: pc,
	}, nil
}

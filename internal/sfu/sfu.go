package sfu

import (
	"github.com/obrel/aira-websocket-stt/pkg/transcribe"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/pion/webrtc/v4"
)

type SFU struct {
	transcriber  transcribe.Service
	trackHandler func(transcribe.Stream, *webrtc.TrackRemote, *webrtc.DataChannel) error
}

func NewSFU(transcriber transcribe.Service, trackHandler func(transcribe.Stream, *webrtc.TrackRemote, *webrtc.DataChannel) error) Service {
	return &SFU{
		transcriber:  transcriber,
		trackHandler: trackHandler,
	}
}

func (s *SFU) CreatePeerConnection() (PeerConnection, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, err
	}

	dataChan := make(chan *webrtc.DataChannel)

	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})
	if err != nil {
		log.For("sfu", "peer").Error(err)
		return nil, err
	}

	pc.OnTrack(func(track *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		log.Printf("Received audio (%s) track, id = %s\n", track.Codec().MimeType, track.ID())

		stream, err := s.transcriber.CreateStream()
		if err != nil {
			log.For("sfu", "peer").Error(err)
		} else {
			err := s.trackHandler(stream, track, <-dataChan)
			if err != nil {
				log.For("sfu", "peer").Error(err)
			}
		}
	})

	pc.OnICEConnectionStateChange(func(connState webrtc.ICEConnectionState) {
		log.Printf("Connection state: %s \n", connState.String())
	})

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dataChan <- dc
	})

	return &Client{
		pc: pc,
	}, nil
}

package sfu

import (
	"github.com/obrel/go-lib/pkg/log"
	"github.com/pion/webrtc/v4"
)

type SFU struct {
	TrackRemote chan *webrtc.TrackRemote
	DataChannel chan *webrtc.DataChannel
}

func NewSFU() *SFU {
	return &SFU{
		TrackRemote: make(chan *webrtc.TrackRemote),
		DataChannel: make(chan *webrtc.DataChannel),
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

		s.TrackRemote <- track
		s.DataChannel = dataChan
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

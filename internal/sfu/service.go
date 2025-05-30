package sfu

import (
	"io"
)

type PeerConnection interface {
	io.Closer
	GetOffer() (string, error)
	SetOffer(offer string) error
	GetAnswer() (string, error)
	SetAnswer(answer string) error
}

type Service interface {
	CreatePeerConnection() (PeerConnection, error)
}

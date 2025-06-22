package handler

import (
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
)

type newSessionRequest struct {
	Offer string `json:"offer"`
}

type newSessionResponse struct {
	Answer string `json:"answer"`
}

type newResultsResponse struct {
	Results []transcription.Result `json:"results"`
}

package handler

import "github.com/obrel/aira-websocket-stt/internal/transcribe"

type newSessionRequest struct {
	Offer string `json:"offer"`
}

type newSessionResponse struct {
	Answer string `json:"answer"`
}

type newResultsResponse struct {
	Results []transcribe.Result `json:"results"`
}

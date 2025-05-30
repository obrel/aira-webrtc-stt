package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/obrel/aira-websocket-stt/internal/sfu"
)

func NewHandler(sfu sfu.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/signaling", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		req := newSessionRequest{}

		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		peer, err := sfu.CreatePeerConnection()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = peer.SetOffer(req.Offer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		answer, err := peer.GetAnswer()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		payload, err := json.Marshal(newSessionResponse{
			Answer: answer,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(payload)
	})

	fs := http.FileServer(http.Dir("view"))
	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}

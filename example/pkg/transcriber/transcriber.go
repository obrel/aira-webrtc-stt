package transcriber

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/obrel/aira-websocket-stt/pkg/decoder"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/pion/webrtc/v4"
)

func Transcribe(trans transcription.Transcription, sampleRate int, track *webrtc.TrackRemote, dc *webrtc.DataChannel) error {
	decoder, err := decoder.NewDecoder(sampleRate)
	if err != nil {
		return err
	}

	err = trans.Connect()
	if err != nil {
		return err
	}

	defer func() {
		err := trans.Close()
		if err != nil {
			log.For("transcriber", "transcribe").Error(err)
			return
		}

		dc.Close()
	}()

	errs := make(chan error, 2)
	audioStream := make(chan []byte)
	response := make(chan bool)
	result := make(chan transcription.Result)
	doneWrite := make(chan bool)
	doneTranscribe := make(chan bool)
	timer := time.NewTimer(5 * time.Second)

	go func() {
		err := trans.Receive(result, doneTranscribe)
		if err != nil {
			log.For("transcriber", "transcribe").Error(err)
			return
		}
	}()

	go func() {
		for {
			packet, _, err := track.ReadRTP()
			timer.Reset(1 * time.Second)

			if err != nil {
				timer.Stop()

				if err == io.EOF {
					doneTranscribe <- true
					doneWrite <- true
					return
				}

				errs <- err
				return
			}

			audioStream <- packet.Payload
			<-response
		}
	}()

	for {
		select {
		case audioChunk := <-audioStream:
			response <- true

			if len(audioChunk) > 0 {
				payload, err := decoder.Decode(audioChunk)
				if err != nil {
					return err
				}

				_, err = trans.Write(payload)
				if err != nil {
					return err
				}
			}
		case result := <-result:
			msg, err := json.Marshal(result)
			if err != nil {
				continue
			}

			if string(msg) != "" {
				err = dc.Send(msg)
				if err != nil {
					log.For("transcriber", "transcribe").Error(err)
				}
			}
		case <-doneWrite:
			return nil
		case <-timer.C:
			return fmt.Errorf("Read operation timed out")
		case err = <-errs:
			log.Printf("Unexpected error reading track %s: %v", track.ID(), err)
			return err
		}
	}
}

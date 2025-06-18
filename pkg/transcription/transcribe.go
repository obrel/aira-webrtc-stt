package transcription

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/obrel/aira-websocket-stt/pkg/decoder"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/pion/webrtc/v4"
)

func Transcribe(transcription Transcription, sampleRate int, track *webrtc.TrackRemote, dc *webrtc.DataChannel) error {
	decoder, err := decoder.NewDecoder(sampleRate)
	if err != nil {
		return err
	}

	defer func() {
		err := transcription.Close()
		if err != nil {
			log.For("transcription", "transcribe").Error(err)
			return
		}

		dc.Close()
	}()

	errs := make(chan error, 2)
	audioStream := make(chan []byte)
	response := make(chan bool)
	result := make(chan Result)
	done := make(chan bool)
	timer := time.NewTimer(5 * time.Second)

	go func() {
		err := transcription.Receive(result, done)
		if err != nil {
			log.For("transcription", "transcribe").Error(err)
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
					done <- true
					close(audioStream)
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

				_, err = transcription.Write(payload)
				if err != nil {
					return err
				}
			}
		case result := <-result:
			msg, err := json.Marshal(result)
			if err != nil {
				continue
			}

			err = dc.Send(msg)
			if err != nil {
				log.For("transcription", "transcribe").Error(err)
			}
		case <-done:
			return nil
		case <-timer.C:
			return fmt.Errorf("Read operation timed out")
		case err = <-errs:
			log.Printf("Unexpected error reading track %s: %v", track.ID(), err)
			return err
		}
	}

	fmt.Println("ASU")
	return nil
}

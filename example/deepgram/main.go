package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/obrel/aira-websocket-stt/example/pkg/handler"
	"github.com/obrel/aira-websocket-stt/example/pkg/sfu"
	"github.com/obrel/aira-websocket-stt/pkg/transcription"
	"github.com/obrel/aira-websocket-stt/pkg/transcription/service/deepgram"
	"github.com/obrel/go-lib/pkg/log"
)

const (
	httpDefaultPort   = "4000"
	defaultStunServer = "stun:stun.l.google.com:19302"
)

var (
	apiKey *string
)

func main() {
	httpPort := flag.String("http-port", httpDefaultPort, "HTTP listen port")
	apiKey = flag.String("api-key", "", "Deepgram api key")
	flag.Parse()

	ctx := context.Background()
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	tr, err := transcription.NewTranscription("deepgram", []transcription.Option{
		deepgram.ApiKey(*apiKey),
		deepgram.Model("nova-2"),
		deepgram.Language("en"),
		deepgram.Encoding("linear16"),
		deepgram.SampleRate(16000),
	}...)
	if err != nil {
		log.For("aira", "main").Fatal(err)
	}

	sfu := sfu.NewSFU(tr)
	server := &http.Server{Addr: fmt.Sprintf(":%v", *httpPort), Handler: handler.NewHandler(sfu)}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, cancelCtx := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancelCtx()

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.For("aira", "main").Info("Graceful shutdown timed out. Forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.For("aira", "main").Error(err)
		}

		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.For("main", "run").Error(err)
		return
	}

	<-serverCtx.Done()
}

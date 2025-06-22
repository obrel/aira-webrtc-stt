# AIRA STT - Cartesia

## How to run
```
$ go run main.go --api-key <cartesia_api_key>
```

Open your browser, and go to http://localhost:4000. Click start, start talking, and you'll be receive the transcription text.

## Snippet
```
import (
    "github.com/obrel/aira-websocket-stt/pkg/transcription"
    "github.com/obrel/aira-websocket-stt/pkg/transcription/service/deepgram"
)

func main() {
    ...

    // initializa transcription
    trans, err := transcription.NewTranscription("cartesia", []transcription.Option{
		cartesia.ApiKey(apiKey),
		cartesia.Model("ink-whisper"),
		cartesia.Language("en"),
		cartesia.Encoding("pcm_s16le"),
		cartesia.SampleRate(16000),
	}...)
	if err != nil {
		log.Fatal(err)
	}

    // connect to transcription service
    err = trans.Connect()
	if err != nil {
		log.Fatal(err)
	}

    // listen to transcription result
    result := make(chan transcription.Result)
	doneTranscribe := make(chan bool, 1)

	go func() {
		err := trans.Receive(result, doneTranscribe)
		if err != nil {
			log.Error(err)
			return
		}
	}()

    // send your audio payload []byte
    _, err = trans.Write(payload)
    if err != nil {
        log.Error(err)
    }

    // print result
    for {
		select {
        case result := <-result:
            msg, err := json.Marshal(result)
            if err != nil {
                log.Error(err)
                continue
            }

            if string(msg) != "" {
                if err != nil {
                    log.Error(err)
                }

                fmt.Println(string(msg))
            }
        }
    }

    ...
}
```
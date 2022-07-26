package handerls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sehogas/sse/util"
)

type Message struct {
	Origin string    `json:"origin,omitempty"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time,omitempty"`
}

func SendEvent(s util.Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go func() {
			message.Time = time.Now()
			b, _ := json.Marshal(message)
			s.Notify(b)
			/*
				if err := s.Notify(b); err != nil {
					log.error(err.Error())
					return
				}
			*/
		}()

		w.WriteHeader(http.StatusCreated)
	}
}

func ServerSentEvent(s util.Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Subscribe
		c := make(chan []byte)
		unsubscribeFn, err := s.Subscribe(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Signal SSE Support
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		//w.Header().Set("Access-Control-Allow-Origin", "*")

	Looping:
		for {
			select {
			case <-r.Context().Done():
				if err := unsubscribeFn(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				break Looping

			default:
				b := <-c
				fmt.Fprintf(w, "data: %s\n\n", b)

				w.(http.Flusher).Flush()
			}
		}
	}
}

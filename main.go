package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sehogas/sse/util"
)

type Mensaje struct {
	Origen string    `bson:"origen" json:"origen,omitempty"`
	Texto  string    `bson:"texto" json:"texto,omitempty"`
	Fecha  time.Time `bson:"fecha" json:"fecha,omitempty"`
}

func main() {
	nc := util.NewNotificationCenter()

	/*
		go func() {
			for {
				b := []byte(time.Now().Format(time.RFC3339))
				if err := nc.Notify(b); err != nil {
					log.Fatal(err)
				}

				time.Sleep(2 * time.Second)
			}
		}()
	*/

	router := mux.NewRouter()
	router.HandleFunc("/sse", handleSSE(nc))
	router.HandleFunc("/sendmessage", handleSendEvent(nc)).Methods("POST")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	handler := cors.AllowAll().Handler(router)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}

func handleSendEvent(s util.Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mensaje Mensaje
		err := json.NewDecoder(r.Body).Decode(&mensaje)
		if err != nil {
			http.Error(w, "parámetros del body incorrectos. "+err.Error(), http.StatusBadRequest)
			return
		}
		go func() {
			mensaje.Fecha = time.Now()
			b, _ := json.Marshal(mensaje)
			s.Notify(b)
			/*
				if err := s.Notify(b); err != nil {
					//http.Error(w, "error enviando mensaje. "+err.Error(), http.StatusInternalServerError)
					log.error("error enviando mensaje. "+err.Error())
					return
				}
			*/
		}()

		w.WriteHeader(http.StatusCreated)
	}
}

func handleSSE(s util.Subscriber) http.HandlerFunc {
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

/*
type Notifier interface {
	Notify(b []byte) error
}
*/

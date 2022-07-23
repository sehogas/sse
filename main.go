package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sehogas/sse_server/handerls"
	"github.com/sehogas/sse_server/util"
)

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
	router.HandleFunc("/sse", handerls.ServerSentEvent(nc))
	router.HandleFunc("/sendmessage", handerls.SendEvent(nc)).Methods("POST")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3003"
	}

	handler := cors.AllowAll().Handler(router)

	log.Println("Listening and serving in port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}

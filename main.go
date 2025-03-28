package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/sehogas/sse/handerls"
	"github.com/sehogas/sse/util"
)

var (
	Version string = "development"
)

type Middleware func(http.Handler) http.Handler

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)
		log.Println(wrapper.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	nc := util.NewNotificationCenter()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3003"
	}

	middlewareCors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		Debug:            false,
	})

	stack := CreateStack(
		middlewareCors.Handler,
		Logging,
	)

	router := http.NewServeMux()
	router.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Version string `json:"version"`
		}{
			Version: Version,
		}
		json, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	})
	router.HandleFunc("/sse", handerls.ServerSentEvent(nc))
	router.HandleFunc("POST /sendmessage", handerls.SendEvent(nc))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", PORT),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go gracefulShutdown(server)

	log.Printf("Starting server on port %v", PORT)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func gracefulShutdown(apiServer *http.Server) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
}

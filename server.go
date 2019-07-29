package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Server struct {
}

func (s *Server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleWebhookVerification(w, r)
			return
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error Reading Body")
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("Message Received. Webhook input.")
		log.Printf("%v", string(b))
		w.Write([]byte("Hello World"))
		w.WriteHeader(http.StatusOK)
	}
}

func handleWebhookVerification(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	token := r.URL.Query().Get("hub.verify_token")

	if token == os.Getenv("VERIFY_TOKEN") {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}

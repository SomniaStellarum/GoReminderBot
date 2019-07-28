package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
}

func (s *Server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error Reading Body")
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("Message Received. Webhook input.")
		log.Printf(string(b))
		w.WriteHeader(http.StatusOK)
	}
}

package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type Server struct {
}

func (s *Server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleWebhookVerification(w, r)
			return
		}
		ctx := appengine.NewContext(r)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(ctx, "Error Reading Body")
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Infof(ctx, "Message Received. Webhook input.")
		log.Infof(ctx, "%v", string(b))
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

package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
)

type Server struct {
	df              *dialogflow.SessionsClient
	dataClient      *datastore.Client
	verifyToken     string
	projectID       string
	pageAccessToken string
	chMessages      chan []byte
}

func newServer(df *dialogflow.SessionsClient) (s *Server, err error) {
	ctx := context.Background()
	s = new(Server)
	s.df = df
	s.verifyToken = os.Getenv("VERIFY_TOKEN")
	s.projectID = os.Getenv("PROJECT_ID")
	s.pageAccessToken = os.Getenv("PAGE_ACCESS_TOKEN")
	s.dataClient, err = datastore.NewClient(ctx, s.projectID)
	if err != nil {
		return nil, err
	}
	s.chMessages = make(chan []byte)
	return s, nil
}

func (s *Server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			s.handleWebhookVerification()(w, r)
			return
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error Reading Body")
			return
		}
		log.Printf("Message Received. Webhook input.")
		log.Printf("%v", string(b))
		//s.runMessaging(b)
		s.chMessages <- b
	}
}

func (s *Server) handleWebhookVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("hub.challenge")
		token := r.URL.Query().Get("hub.verify_token")
		if token == s.verifyToken {
			w.WriteHeader(200)
			w.Write([]byte(challenge))
		} else {
			w.WriteHeader(404)
			w.Write([]byte("Error, wrong validation token"))
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
)

func main() {
	ctx := context.Background()
	cl, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		log.Fatal("Couldn't open dialogflow client.")
	}
	defer cl.Close()
	s := newServer(cl)
	s.routes()

	go s.runMessaging()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

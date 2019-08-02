package main

import "net/http"

func (s *Server) routes() {
	http.HandleFunc("/webhook", s.handleWebhook())
	http.HandleFunc("/alert", s.handleAlerts())
}

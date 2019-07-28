package main

import "net/http"

func (s *Server) routes() {
	http.HandleFunc("/webhoook", s.handleWebhook())
}

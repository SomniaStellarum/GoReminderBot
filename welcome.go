package main

func (s *Server) sendWelcome(id string) {
	m := NewWelcome(id)
	s.sendAlertMessage(m)
}

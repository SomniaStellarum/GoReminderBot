package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) runMessaging(b []byte) {
	// Goroutine to handle all messaging
	// Webhook receives message and passes it to this goroutine before
	// returning status 200
	//for {
	//b := <-s.chMessages
	obj := new(Object)
	err := json.Unmarshal(b, obj)
	if err != nil {
		log.Printf("Error Unmarshalling Message: %v", err)
		//continue
		return
	}
	for _, e := range obj.Entries {
		for _, m := range e.Messages {
			sender := m.Sender.ID
			//recipient := m.Recipient.ID
			text := m.Mes.Text
			reply, err := s.parseText(sender, text)
			if err != nil {
				log.Printf("Error Parsing Text: %v", err)
				continue
			}
			m := NewResponseMessage(sender, reply)
			s.sendMessage(m)
		}
	}
	//}
}

func (s *Server) sendMessage(m *MessageBody) {
	log.Printf("Sending Message")
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error Marshalling Reply: %v", err)
		return
	}
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", "https://graph.facebook.com/v4.0/me/messages", buf)
	if err != nil {
		log.Printf("Error Creating Reply Request: %v", err)
		return
	}
	q := req.URL.Query()
	q.Add("access_token", s.pageAccessToken)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error Replying: %v\n%v", err, res)
		return
	}
	log.Printf("Message Sent")
	log.Printf("Response: %v", res)
}

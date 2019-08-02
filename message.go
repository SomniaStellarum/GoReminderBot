package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) runMessaging() {
	// Goroutine to handle all messaging
	// Webhook receives message and passes it to this goroutine before
	// returning status 200
	for {
		select {
		case b := <-s.chMessages:
			obj := new(Object)
			err := json.Unmarshal(b, obj)
			if err != nil {
				log.Printf("Error Unmarshalling Message: %v", err)
				continue
			}
			for _, e := range obj.Entries {
				for _, m := range e.Messages {
					if m.Mes.Text != "" {
						sender := m.Sender.ID
						text := m.Mes.Text
						s.manageMessage(sender, text)
					} else if m.Pb.P != "" {
						sender := m.Sender.ID
						text := m.Pb.P
						s.manageMessage(sender, text)
					} else {
						log.Printf("Message format error. No Postback or Message")
					}
				}
			}
		}
	}
}

func (s *Server) manageMessage(sender, text string) {
	if !s.isUser(sender) {
		s.sendWelcome(sender)
	}
	reply, queryResult, err := s.parseText(sender, text)
	if err != nil {
		log.Printf("Error Parsing Text: %v", err)
		return
	}
	if queryResult.GetAllRequiredParamsPresent() {
		switch action := queryResult.GetAction(); action {
		case "reminders.add":
			var t time.Time
			var name string
			params := queryResult.GetParameters()
			fields := params.GetFields()
			for f, v := range fields {
				switch f {
				case "date-time":
					s := v.GetStringValue()
					t, err = time.Parse(time.RFC3339, s)
					if err != nil {
						log.Printf("Error Formating Time: %v", err)
					}

					location, err := time.LoadLocation("Local")
					if err != nil {
						log.Printf("Error Getting Local Time")
					}
					t = t.In(location)
				case "name":
					name = v.GetStringValue()
				}
			}
			err = s.storeReminder(sender, name, t)
			if err != nil {
				log.Printf("Error Storing Data: %v", err)
			}
			m := NewResponseMessage(sender, reply)
			s.sendMessage(m)
		case "reminders.get":
			rem, err := s.getReminders(sender)
			if err != nil {
				log.Printf("Error getting Reminders: %v", err)
				return
			}
			for _, r := range rem {
				log.Printf("Reminder: %v", r.String())
				m := NewResponseMessage(sender, r.String())
				s.sendMessage(m)
			}
		case "reminders.snooze":
			rem, err := s.getReminders(sender)
			if err != nil {
				log.Printf("Error getting Reminders: %v", err)
				return
			}
			for _, r := range rem {
				if r.Status == StatusAlert {
					r.Status = StatusSnoozed
					params := queryResult.GetParameters()
					fields := params.GetFields()
					s := fields["date-time"].GetStringValue()
					r.NextTime, err = time.Parse(time.RFC3339, s)
					if err != nil {
						log.Printf("Error Parsing Time: %v", err)
					}
				}
			}
		case "reminders.remove":
			rem := make([]*Reminder, 0)
			s.updateReminders(sender, rem)
			m := NewResponseMessage(sender, reply)
			s.sendMessage(m)
		case "reminders.accept":
			rem, err := s.getReminders(sender)
			if err != nil {
				log.Printf("Error getting Reminders: %v", err)
				return
			}
			for _, r := range rem {
				if r.Status == StatusAlert {
					r.Status = StatusDone
					params := queryResult.GetParameters()
					fields := params.GetFields()
					s := fields["date-time"].GetStringValue()
					r.NextTime, err = time.Parse(time.RFC3339, s)
					if err != nil {
						log.Printf("Error Parsing Time: %v", err)
					}
				}
			}
		default:
			m := NewResponseMessage(sender, reply)
			s.sendMessage(m)
		}
	} else {
		m := NewResponseMessage(sender, reply)
		s.sendMessage(m)
	}
}

func (s *Server) sendMessage(m *MessageBody) {
	printMessage(m)
	if !s.debugMode {
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
}

func printMessage(m *MessageBody) {
	log.Printf("Message: %v\nTo Recipient: %v", m.Mes, m.Recipient)
}

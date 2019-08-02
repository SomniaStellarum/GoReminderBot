package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) sendAlerts() {
	data, err := s.getAllReminders()
	if err != nil {
		log.Printf("Error Getting Data for Alerts: %v", err)
	}
	t := time.Now()
	for _, d := range data {
		log.Printf("Checking Reminders for %v. Number: %v", d.UserID, len(d.Reminders))
		var rems []string
		alert := false
		for _, r := range d.Reminders {
			preAlarm := (r.Status == StatusPre) && (r.SetTime.After(t))
			snoozeAlarm := (r.Status == StatusSnoozed) && (r.NextTime.After(t))
			log.Printf("Reminder: %v Pre?: %v Snooze?: %v", r.Desc, preAlarm, snoozeAlarm)
			if preAlarm || snoozeAlarm {
				alert = true
				rems = append(rems, r.Desc)
				r.Status = StatusAlert
			}
		}
		if alert {
			// Send Alert
			m := NewMessageAttach(d.UserID, rems)
			s.sendAlertMessage(m)
			// Update datastore
			s.updateReminders(d.UserID, d.Reminders)
		}
	}
}

func (s *Server) sendAlertMessage(m *MessageAttach) {
	log.Printf("Alert Sending")
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

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	StatusPre int = iota + 1
	StatusAlert
	StatusSnoozed
	StatusDone
)

type Reminder struct {
	Desc     string    `datastore:"description"`
	SetTime  time.Time `datastore:"set_time"`
	Status   int       `datastore:"status"`
	NextTime time.Time `datastore:"next_time"`
}

type UserData struct {
	UserID    *datastore.Key `datastore:"__key__"`
	Reminders []*Reminder    `datastore:"reminders"`
}

type Data struct {
	UserID    string
	Reminders []*Reminder
}

func (s *Server) storeReminder(id, name string, t time.Time) error {
	log.Printf("Storing Reminder. \nName: %v\nTime: %v", name, t)
	// Setup Reminder
	rem := new(Reminder)
	rem.Desc = name
	rem.SetTime = t
	rem.Status = StatusPre

	// Store in Datastore. Check ID exists first
	ctx := context.Background()
	data := new(UserData)
	k := datastore.NameKey("Reminder", id, nil)
	err := s.dataClient.Get(ctx, k, data)
	if err == datastore.ErrNoSuchEntity {
		data.UserID = k
	} else if err != nil {
		log.Printf("Error Getting Data from Datastore: %v", err)
		return err
	}
	data.Reminders = append(data.Reminders, rem)
	_, err = s.dataClient.Put(ctx, k, data)
	if err != nil {
		log.Printf("Error Storing Data in Datastore: %v", err)
	}
	return err
}

func (s *Server) getReminders(id string) ([]*Reminder, error) {
	ctx := context.Background()
	data := new(UserData)
	k := datastore.NameKey("Reminder", id, nil)
	err := s.dataClient.Get(ctx, k, data)
	if err != nil {
		log.Printf("Error Getting Data: %v", err)
		return nil, err
	}
	return data.Reminders, nil
}

func (s *Server) updateReminders(id string, reminders []*Reminder) error {
	ctx := context.Background()
	data := new(UserData)
	k := datastore.NameKey("Reminder", id, nil)
	data.Reminders = reminders
	data.UserID = k
	_, err := s.dataClient.Put(ctx, k, data)
	return err
}

func (s *Server) isUser(id string) bool {
	ctx := context.Background()
	data := new(UserData)
	k := datastore.NameKey("Reminder", id, nil)
	err := s.dataClient.Get(ctx, k, data)
	if err == datastore.ErrNoSuchEntity {
		return false
	}
	return true
}

func (s *Server) getAllReminders() ([]*Data, error) {
	ctx := context.Background()
	q := datastore.NewQuery("")
	data := make([]*UserData, 0)
	_, err := s.dataClient.GetAll(ctx, q, data)
	if err != nil {
		return nil, err
	}
	outData := make([]*Data, 0)
	for _, d := range data {
		outD := new(Data)
		outD.Reminders = d.Reminders
		outD.UserID = d.UserID.String()
	}
	return outData, nil
}

func (r *Reminder) String() string {
	var status string
	switch r.Status {
	case StatusPre:
		status = "Set"
	case StatusDone:
		status = "Complete"
	case StatusAlert:
		status = "Ringing"
	case StatusSnoozed:
		status = "Snoozed"
	}
	return fmt.Sprintf(
		"Reminder: %v --- Time: %v --- Status: %v",
		r.Desc,
		r.SetTime.Format("Mon Jan _2 - 3:04PM"),
		status,
	)
}

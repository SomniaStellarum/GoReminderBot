package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	StatusPre int = iota + 1
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

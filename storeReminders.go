package main

import (
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
	Reminders []Reminder     `datastore:"reminders"`
}

func (s *Server) storeReminder(id, name string, t time.Time) error {
	log.Printf("Storing Reminder. \nName: %v\nTime: %v", name, t)
	return nil
}

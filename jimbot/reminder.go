package jimbot

import (
	"log"
	"time"
)

// Reminder type for reminders data
type Reminder struct {
	notifyTime time.Time
	event      string
}

// NewReminder : create new reminder with time and event
func NewReminder(notifyTime time.Time, event string) Reminder {
	var reminder Reminder
	reminder.notifyTime = notifyTime
	reminder.event = event
	log.Print("[+++] REMINDER ", notifyTime.UTC(), event)
	return reminder
}

// Notifier : wait until it's time to send our reminder
func Notifier(reminder Reminder) string {
	for {
		if time.Now().UTC() == reminder.notifyTime.UTC() {
			break
		}
		time.Sleep(1 * time.Minute)
	}
	notification := reminder.event
	return notification
}

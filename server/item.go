package main

import (
	"fmt"
	"time"
)

const (
	// StatusEnqueued denotes items on the own list waiting to be done (empty string for legacy)
	StatusEnqueued = ""
	// StatusComplete denotes items sent that are finished by the receiver
	StatusComplete = "Complete"
	// StatusDeleted denotes items sent that are deleted by the receiver
	StatusDeleted = "Deleted"
	// StatusPending denotes items sent that are not yet processed by the receiver
	StatusPending = "Pending"
)

// Item represents a to do item
type Item struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	CreateAt int64  `json:"create_at"`
	CreateBy string `json:"create_by"`
	SendTo   string `json:"send_to"`
	Status   string `json:"status"`
}

func itemsListToString(items []*Item) string {
	if len(items) == 0 {
		return "Nothing to do!"
	}

	str := "\n\n"

	for _, item := range items {
		createAt := time.Unix(item.CreateAt/1000, 0)
		str += fmt.Sprintf("* %s\n  * (%s)\n", item.Message, createAt.Format("January 2, 2006 at 15:04"))
	}

	return str
}

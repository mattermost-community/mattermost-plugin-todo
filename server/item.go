package main

import (
	"fmt"
	"time"
)

// Item represents a to do item
type Item struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	CreateAt int64  `json:"create_at"`
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

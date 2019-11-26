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

func (p *Plugin) addItem(userID string, item *Item) error {
	err := p.storeItem(item)
	if err != nil {
		return err
	}

	p.addToOrderForUser(userID, item.ID)
	if err != nil {
		p.deleteItem(item.ID)
		return err
	}

	return nil
}

func itemsListToString(items []*Item) string {
	str := "To Do List:\n\n"

	if len(items) == 0 {
		return str + "Nothing to do!"
	}

	for _, item := range items {
		createAt := time.Unix(item.CreateAt/1000, 0)
		str += fmt.Sprintf("* %s (added %s)\n", item.Message, createAt.Format("January 2, 2006"))
	}

	return str
}

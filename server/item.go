package main

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

// Item represents a to do item
type Item struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	CreateAt int64  `json:"create_at"`
}

// ExtendedItem extends the information on Item to be used on the front-end
type ExtendedItem struct {
	Item
	ForeignUser     string `json:"user"`
	ForeignList     string `json:"list"`
	ForeignPosition int    `json:"position"`
}

func newItem(message string) *Item {
	return &Item{
		ID:       model.NewId(),
		CreateAt: model.GetMillis(),
		Message:  message,
	}
}

func itemsListToString(items []*ExtendedItem) string {
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

func (p *Plugin) extendItemInfo(item *Item, oe *OrderElement) *ExtendedItem {
	if item == nil || oe == nil {
		return nil
	}

	feItem := &ExtendedItem{
		Item: *item,
	}

	if oe.ForeignUserID == "" {
		return feItem
	}

	_, _, n, listKey := p.getUserListForItem(oe.ForeignUserID, oe.ForeignItemID)

	var listName string
	switch listKey {
	case MyListKey:
		listName = ""
	case InListKey:
		listName = "in"
	case OutListKey:
		listName = "out"
	}

	userName := "Someone"
	if user, err := p.API.GetUser(oe.ForeignUserID); err == nil {
		userName = user.Username
	}

	feItem.ForeignUser = userName
	feItem.ForeignList = listName
	feItem.ForeignPosition = n

	return feItem
}

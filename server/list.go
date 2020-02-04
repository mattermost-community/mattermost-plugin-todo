package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	// MyListKey is the key used to store the order of the owned todos
	MyListKey = ""
	// InListKey is the key used to store the order of received todos
	InListKey = "_in"
	// OutListKey is the key used to store the order of sent todos
	OutListKey = "_out"
)

// ListStore represents the KVStore operations for lists
type ListStore interface {
	// Item related function
	AddItem(item *Item) error
	GetItem(itemID string) (*Item, error)
	RemoveItem(itemID string) error

	// GetItemOrder gets the item Order Element and position of the item itemID on user userID's list listID
	GetItemOrder(userID, itemID, listID string) (*OrderElement, int, error)
	// GetItemListAndOrder gets the item list, Order Element and position for user userID
	GetItemListAndOrder(userID, itemID string) (string, *OrderElement, int)

	// Order Element related functions

	// Add creates a new OrderElement with the itemID, foreignUSerID and foreignItemID, and stores it
	// on the listID for userID.
	Add(userID, itemID, listID, foreignUserID, foreignItemID string) error
	// Remove removes the OrderElement for itemID in listID for userID
	Remove(userID, itemID, listID string) error
	// Pop removes the first OrderElement in listID for userID and return it
	Pop(userID, listID string) (*OrderElement, error)

	// GetList returns the list of OrderElement in listID for userID
	GetList(userID, listID string) ([]*OrderElement, error)
}

type listManager struct {
	store ListStore
	api   plugin.API
}

// NewListManager creates a new listManager
func NewListManager(store ListStore, api plugin.API) *listManager {
	return &listManager{
		store: store,
		api:   api,
	}
}

func (l *listManager) Add(userID, message string) error {
	item := newItem(message)

	if err := l.store.AddItem(item); err != nil {
		return err
	}

	if err := l.store.Add(userID, item.ID, MyListKey, "", ""); err != nil {
		l.store.RemoveItem(item.ID)
		return err
	}

	return nil
}

func (l *listManager) Send(senderID, receiverID, message string) (string, error) {
	senderItem := newItem(message)
	if err := l.store.AddItem(senderItem); err != nil {
		return "", err
	}

	receiverItem := newItem(message)
	if err := l.store.AddItem(receiverItem); err != nil {
		l.store.RemoveItem(senderItem.ID)
		return "", err
	}

	appErr := l.store.Add(senderID, senderItem.ID, OutListKey, receiverID, receiverItem.ID)
	if appErr != nil {
		l.store.RemoveItem(senderItem.ID)
		l.store.RemoveItem(receiverItem.ID)
		return "", appErr
	}

	appErr = l.store.Add(receiverID, receiverItem.ID, InListKey, senderID, senderItem.ID)
	if appErr != nil {
		l.store.RemoveItem(senderItem.ID)
		l.store.RemoveItem(receiverItem.ID)
		l.store.Remove(senderID, senderItem.ID, OutListKey)
		return "", appErr
	}

	return receiverItem.ID, nil
}

func (l *listManager) Get(userID, listID string) ([]*ExtendedItem, error) {
	oes, err := l.store.GetList(userID, listID)
	if err != nil {
		return nil, err
	}

	extendedItems := []*ExtendedItem{}
	for _, oe := range oes {
		item, err := l.store.GetItem(oe.ItemID)
		if err != nil {
			continue
		}

		extendedItem := l.extendItemInfo(item, oe)
		extendedItems = append(extendedItems, extendedItem)
	}

	return extendedItems, nil
}

func (l *listManager) Complete(userID, itemID string) (todoMessage string, foreignUserID string, outErr error) {
	itemList, oe, _ := l.store.GetItemListAndOrder(userID, itemID)
	if oe == nil {
		return "", "", fmt.Errorf("cannot find element")
	}

	if err := l.store.Remove(userID, itemID, itemList); err != nil {
		return "", "", err
	}

	l.store.RemoveItem(itemID)

	if oe.ForeignUserID == "" {
		return "", "", nil
	}

	l.store.Remove(oe.ForeignUserID, oe.ForeignItemID, OutListKey)
	item, err := l.store.GetItem(oe.ForeignItemID)
	if err != nil {
		return "", "", nil
	}

	l.store.RemoveItem(oe.ForeignItemID)

	return item.Message, oe.ForeignUserID, nil
}

func (l *listManager) Enqueue(userID, itemID string) (todoMessage string, foreignUserID string, outErr error) {
	item, err := l.store.GetItem(itemID)
	if err != nil {
		return "", "", err
	}

	oe, _, err := l.store.GetItemOrder(userID, itemID, InListKey)
	if err != nil {
		return "", "", err
	}
	if oe == nil {
		return "", "", fmt.Errorf("element order not found")
	}

	err = l.store.Add(userID, itemID, MyListKey, oe.ForeignUserID, oe.ForeignItemID)
	if err != nil {
		return "", "", err
	}

	err = l.store.Remove(userID, itemID, InListKey)
	if err != nil {
		l.store.Remove(userID, itemID, MyListKey)
		return "", "", err
	}

	return item.Message, oe.ForeignUserID, nil
}

func (l *listManager) Remove(userID, itemID string) (todoMessage string, foreignUserID string, isSender bool, outErr error) {
	itemList, oe, _ := l.store.GetItemListAndOrder(userID, itemID)
	if oe == nil {
		return "", "", false, fmt.Errorf("cannot find element")
	}

	if err := l.store.Remove(userID, itemID, itemList); err != nil {
		return "", "", false, err
	}

	l.store.RemoveItem(itemID)

	if oe.ForeignUserID == "" {
		return "", "", false, nil
	}

	list, _, _ := l.store.GetItemListAndOrder(oe.ForeignUserID, oe.ForeignItemID)

	l.store.Remove(oe.ForeignUserID, oe.ForeignItemID, list)
	item, err := l.store.GetItem(oe.ForeignItemID)
	if err != nil {
		return "", "", false, nil
	}

	l.store.RemoveItem(oe.ForeignItemID)

	return item.Message, oe.ForeignUserID, list == OutListKey, nil
}

func (l *listManager) Pop(userID string) (todoMessage string, sender string, outErr error) {
	oe, err := l.store.Pop(userID, MyListKey)
	if err != nil {
		return "", "", err
	}

	if oe == nil {
		return "", "", nil
	}

	l.store.RemoveItem(oe.ItemID)

	if oe.ForeignUserID == "" {
		return "", "", nil
	}

	item, err := l.store.GetItem(oe.ForeignUserID)
	if err != nil {
		return "", "", nil
	}

	l.store.Remove(oe.ForeignUserID, oe.ForeignItemID, OutListKey)
	l.store.RemoveItem(oe.ForeignItemID)

	return item.Message, oe.ForeignUserID, nil
}

func (l *listManager) GetUserName(userID string) string {
	user, err := l.api.GetUser(userID)
	if err != nil {
		return "Someone"
	}
	return user.Username
}

func (l *listManager) extendItemInfo(item *Item, oe *OrderElement) *ExtendedItem {
	if item == nil || oe == nil {
		return nil
	}

	feItem := &ExtendedItem{
		Item: *item,
	}

	if oe.ForeignUserID == "" {
		return feItem
	}

	list, _, n := l.store.GetItemListAndOrder(oe.ForeignUserID, oe.ForeignItemID)

	var listName string
	switch list {
	case MyListKey:
		listName = ""
	case InListKey:
		listName = "in"
	case OutListKey:
		listName = "out"
	}

	userName := l.GetUserName(oe.ForeignUserID)

	feItem.ForeignUser = userName
	feItem.ForeignList = listName
	feItem.ForeignPosition = n

	return feItem
}

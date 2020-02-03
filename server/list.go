package main

import (
	"fmt"
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
	AddItem(item *Item) error
	GetItem(itemID string) (*Item, error)
	RemoveItem(itemID string) error

	GetItemOrder(userID string, itemID string, listID string) (*OrderElement, int, error)
	GetItemListAndOrder(userID string, itemID string) (string, *OrderElement, int)

	Add(userID string, itemID string, listID string, foreignUserID string, foreignItemID string) error
	Remove(userID string, itemID string, listID string) error
	Pop(userID string, listID string) (*OrderElement, error)

	GetList(userID string, listID string) ([]*OrderElement, error)
}

type listManager struct {
	store       ListStore
	getUserName func(string) string
}

// NewListManager creates a new listManager
func NewListManager(store ListStore, getUserName func(string) string) *listManager {
	return &listManager{
		store:       store,
		getUserName: getUserName,
	}
}

func (l *listManager) Add(userID string, message string) error {
	item := newItem(message)

	l.store.AddItem(item)

	if err := l.store.Add(userID, item.ID, MyListKey, "", ""); err != nil {
		return err
	}

	return nil
}

func (l *listManager) Send(senderID string, receiverID string, message string) error {
	senderItem := newItem(message)
	l.store.AddItem(senderItem)
	receiverItem := newItem(message)
	l.store.AddItem(receiverItem)

	appErr := l.store.Add(senderID, senderItem.ID, OutListKey, receiverID, receiverItem.ID)
	if appErr != nil {
		return appErr
	}

	appErr = l.store.Add(receiverID, receiverItem.ID, InListKey, senderID, senderItem.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (l *listManager) Get(userID string, listID string) ([]*ExtendedItem, error) {
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

func (l *listManager) Complete(userID string, itemID string) (todoMessage string, foreignUserID string, err error) {
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

func (l *listManager) Enqueue(userID string, itemID string) (todoMessage string, foreignUserID string, err error) {
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

func (l *listManager) Remove(userID string, itemID string) (todoMessage string, foreignUserID string, isSender bool, err error) {
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

func (l *listManager) Pop(userID string) (todoMessage string, sender string, err error) {
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

	userName := l.getUserName(oe.ForeignUserID)

	feItem.ForeignUser = userName
	feItem.ForeignList = listName
	feItem.ForeignPosition = n

	return feItem
}

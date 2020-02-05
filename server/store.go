package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	// StoreRetries is the number of retries to use when storing orders fails on a race
	StoreRetries = 3
	// StoreOrderKey is the key used to store orders in the plugin KV store
	StoreOrderKey = "order"
	// StoreItemKey is the key used to store items in the plugin KV store
	StoreItemKey = "item"
	// StoreReminderKey is the key used to store the last time a user was reminded
	StoreReminderKey = "reminder"
)

// OrderElement denotes every element in any of the lists. Contains the item that refers to,
// and may contain foreign ids of item and user, denoting the user this element is related to
// and the item on that user system.
type OrderElement struct {
	ItemID        string `json:"item_id"`
	ForeignItemID string `json:"foreign_item_id"`
	ForeignUserID string `json:"foreign_user_id"`
}

func getListKey(userID string, listID string) string {
	return fmt.Sprintf("%s_%s%s", StoreOrderKey, userID, listID)
}

func getItemKey(itemID string) string {
	return fmt.Sprintf("%s_%s", StoreItemKey, itemID)
}

func getReminderKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreReminderKey, userID)
}

type listStore struct {
	api plugin.API
}

// NewListStore creates a new listStore
func NewListStore(api plugin.API) *listStore {
	return &listStore{
		api: api,
	}
}

func (l *listStore) AddItem(item *Item) error {
	jsonItem, jsonErr := json.Marshal(item)
	if jsonErr != nil {
		return jsonErr
	}

	appErr := l.api.KVSet(getItemKey(item.ID), jsonItem)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (l *listStore) GetItem(itemID string) (*Item, error) {
	originalJSONItem, appErr := l.api.KVGet(getItemKey(itemID))
	if appErr != nil {
		return nil, errors.New(appErr.Error())
	}

	if originalJSONItem == nil {
		return nil, errors.New("cannot find item")
	}

	var item *Item
	err := json.Unmarshal(originalJSONItem, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (l *listStore) RemoveItem(itemID string) error {
	appErr := l.api.KVDelete(getItemKey(itemID))
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (l *listStore) GetItemOrder(userID, itemID, listID string) (*OrderElement, int, error) {
	originalJSONOrder, err := l.api.KVGet(getListKey(userID, listID))
	if err != nil {
		return nil, 0, err
	}

	if originalJSONOrder == nil {
		return nil, 0, errors.New("cannot load list")
	}

	var order []*OrderElement
	jsonErr := json.Unmarshal(originalJSONOrder, &order)
	if jsonErr != nil {
		order, _, jsonErr = l.legacyOrderElement(userID, listID)
		if order == nil {
			return nil, 0, jsonErr
		}
	}

	for i, oe := range order {
		if oe.ItemID == itemID {
			return oe, i, nil
		}
	}
	return nil, 0, errors.New("cannot find item")
}

func (l *listStore) GetItemListAndOrder(userID, itemID string) (string, *OrderElement, int) {
	oe, n, _ := l.GetItemOrder(userID, itemID, MyListKey)
	if oe != nil {
		return MyListKey, oe, n
	}

	oe, n, _ = l.GetItemOrder(userID, itemID, OutListKey)
	if oe != nil {
		return OutListKey, oe, n
	}

	oe, n, _ = l.GetItemOrder(userID, itemID, InListKey)
	if oe != nil {
		return InListKey, oe, n
	}

	return "", nil, 0
}

func (l *listStore) Add(userID, itemID, listID, foreignUserID, foreignItemID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := l.getList(userID, listID)
		if err != nil {
			return err
		}

		for _, oe := range order {
			if oe.ItemID == itemID {
				return errors.New("item id already exists in order")
			}
		}

		order = append(order, &OrderElement{
			ItemID:        itemID,
			ForeignItemID: foreignItemID,
			ForeignUserID: foreignUserID,
		})

		ok, err := l.saveList(userID, listID, order, originalJSONOrder)
		if err != nil {
			return err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return nil
		}
	}

	return errors.New("unable to store installation")
}

func (l *listStore) Remove(userID, itemID, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := l.getList(userID, listID)
		if err != nil {
			return err
		}

		found := false
		for i, oe := range order {
			if oe.ItemID == itemID {
				order = append(order[:i], order[i+1:]...)
				found = true
			}
		}

		if !found {
			return errors.New("cannot find item")
		}

		ok, err := l.saveList(userID, listID, order, originalJSONOrder)
		if err != nil {
			return err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return nil
		}
	}

	return errors.New("unable to store order")
}

func (l *listStore) Pop(userID, listID string) (*OrderElement, error) {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := l.getList(userID, listID)
		if err != nil {
			return nil, err
		}

		if len(order) == 0 {
			return nil, errors.New("cannot find item")
		}

		oe := order[0]
		order = order[1:]

		ok, err := l.saveList(userID, listID, order, originalJSONOrder)
		if err != nil {
			return nil, err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return oe, nil
		}
	}

	return nil, errors.New("unable to store order")
}

func (l *listStore) Bump(userID, itemID, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := l.getList(userID, listID)
		if err != nil {
			return err
		}

		var i int
		var oe *OrderElement

		for i, oe = range order {
			if itemID == oe.ItemID {
				break
			}
		}

		if i == len(order) {
			return errors.New("cannot find item")
		}

		newOrder := append([]*OrderElement{oe}, order[:i]...)
		newOrder = append(newOrder, order[i+1:]...)

		ok, err := l.saveList(userID, listID, newOrder, originalJSONOrder)
		if err != nil {
			return err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return nil
		}
	}

	return errors.New("unable to store order")
}

func (l *listStore) GetList(userID, listID string) ([]*OrderElement, error) {
	oes, _, err := l.getList(userID, listID)
	return oes, err
}

func (l *listStore) getList(userID, listID string) ([]*OrderElement, []byte, error) {
	originalJSONOrder, err := l.api.KVGet(getListKey(userID, listID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONOrder == nil {
		return []*OrderElement{}, originalJSONOrder, nil
	}

	var order []*OrderElement
	jsonErr := json.Unmarshal(originalJSONOrder, &order)
	if jsonErr != nil {
		return l.legacyOrderElement(userID, listID)
	}

	return order, originalJSONOrder, nil
}

func (l *listStore) saveList(userID, listID string, order []*OrderElement, originalJSONOrder []byte) (bool, error) {
	newJSONOrder, jsonErr := json.Marshal(order)
	if jsonErr != nil {
		return false, jsonErr
	}

	ok, appErr := l.api.KVCompareAndSet(getListKey(userID, listID), originalJSONOrder, newJSONOrder)
	if appErr != nil {
		return false, errors.New(appErr.Error())
	}

	return ok, nil
}

func (l *listStore) legacyOrderElement(userID, listID string) ([]*OrderElement, []byte, error) {
	originalJSONOrder, err := l.api.KVGet(getListKey(userID, listID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONOrder == nil {
		return []*OrderElement{}, originalJSONOrder, nil
	}

	var order []string
	jsonErr := json.Unmarshal(originalJSONOrder, &order)
	if jsonErr != nil {
		return nil, nil, jsonErr
	}

	newOrder := []*OrderElement{}
	for _, v := range order {
		newOrder = append(newOrder, &OrderElement{ItemID: v})
	}

	return newOrder, originalJSONOrder, nil
}

func (p *Plugin) saveLastReminderTimeForUser(userID string) error {
	strTime := strconv.FormatInt(model.GetMillis(), 10)
	appErr := p.API.KVSet(getReminderKey(userID), []byte(strTime))
	if appErr != nil {
		return errors.New(appErr.Error())
	}
	return nil
}

func (p *Plugin) getLastReminderTimeForUser(userID string) (int64, error) {
	timeBytes, appErr := p.API.KVGet(getReminderKey(userID))
	if appErr != nil {
		return 0, errors.New(appErr.Error())
	}

	if timeBytes == nil {
		return 0, nil
	}

	reminderAt, err := strconv.ParseInt(string(timeBytes), 10, 64)
	if err != nil {
		return 0, err
	}

	return reminderAt, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
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
	// MyListKey is the key used to store the order of the owned todos
	MyListKey = ""
	// InListKey is the key used to store the order of received todos
	InListKey = "_in"
	// OutListKey is the key used to store the order of sent todos
	OutListKey = "_out"
)

func (p *Plugin) storeLastReminderTimeForUser(userID string) error {
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

func (p *Plugin) storeOutItemForUsers(sender string, receiver string, item *Item) error {
	err := p.storeItem(item)
	if err != nil {
		return err
	}

	err = p.addToInForUser(receiver, item.ID)
	if err != nil {
		p.deleteItem(item.ID)
		return err
	}

	err = p.addToOutForUser(sender, item.ID)
	if err != nil {
		p.removeFromInForUser(receiver, item.ID)
		p.deleteItem(item.ID)
		return err
	}

	return nil
}

func (p *Plugin) storeItemForUser(userID string, item *Item) error {
	err := p.storeItem(item)
	if err != nil {
		return err
	}

	err = p.addToOrderForUser(userID, item.ID)
	if err != nil {
		p.deleteItem(item.ID)
		return err
	}

	return nil
}

func (p *Plugin) storeItem(item *Item) error {
	jsonItem, jsonErr := json.Marshal(item)
	if jsonErr != nil {
		return jsonErr
	}

	appErr := p.API.KVSet(getItemKey(item.ID), jsonItem)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (p *Plugin) getItem(itemID string) (*Item, []byte, error) {
	originalJSONItem, appErr := p.API.KVGet(getItemKey(itemID))
	if appErr != nil {
		return nil, nil, errors.New(appErr.Error())
	}

	if originalJSONItem == nil {
		return nil, nil, nil
	}

	var item *Item
	err := json.Unmarshal(originalJSONItem, &item)
	if err != nil {
		return nil, nil, err
	}

	return item, originalJSONItem, nil
}

func (p *Plugin) deleteItem(itemID string) error {
	appErr := p.API.KVDelete(getItemKey(itemID))
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (p *Plugin) getItemsForUser(userID string) ([]*Item, error) {
	return p.getItemListForUser(userID, MyListKey)
}

func (p *Plugin) getItemListForUser(userID string, listID string) ([]*Item, error) {
	order, _, err := p.getListForUser(userID, listID)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for _, id := range order {
		item, _, err := p.getItem(id)
		if err != nil {
			return nil, err
		}
		if item != nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func (p *Plugin) addToOrderForUser(userID string, itemID string) error {
	return p.addToListForUser(userID, itemID, MyListKey)
}

func (p *Plugin) addToInForUser(userID string, itemID string) error {
	return p.addToListForUser(userID, itemID, InListKey)
}

func (p *Plugin) addToOutForUser(userID string, itemID string) error {
	return p.addToListForUser(userID, itemID, OutListKey)
}

func (p *Plugin) addToListForUser(userID string, itemID string, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getListForUser(userID, listID)
		if err != nil {
			return err
		}

		for _, id := range order {
			if id == itemID {
				return errors.New("item id already exists in order")
			}
		}

		order = append(order, itemID)

		ok, err := p.storeList(userID, listID, order, originalJSONOrder)
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

func (p *Plugin) removeFromOrderForUser(userID string, itemID string) error {
	return p.removeFromListForUser(userID, itemID, MyListKey)
}

func (p *Plugin) removeFromInForUser(userID string, itemID string) error {
	return p.removeFromListForUser(userID, itemID, InListKey)
}

func (p *Plugin) removeFromOutForUser(userID string, itemID string) error {
	return p.removeFromListForUser(userID, itemID, OutListKey)
}

func (p *Plugin) removeFromListForUser(userID string, itemID string, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getListForUser(userID, listID)
		if err != nil {
			return err
		}

		found := false
		for i, id := range order {
			if id == itemID {
				order = append(order[:i], order[i+1:]...)
				found = true
			}
		}

		if !found {
			return nil
		}

		ok, err := p.storeList(userID, listID, order, originalJSONOrder)
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

func (p *Plugin) popFromOrderForUser(userID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getOrderForUser(userID)
		if err != nil {
			return err
		}

		if len(order) == 0 {
			return nil
		}

		order = order[1:]

		ok, err := p.storeOrder(userID, order, originalJSONOrder)
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

func (p *Plugin) storeOrder(userID string, order []string, originalJSONOrder []byte) (bool, error) {
	return p.storeList(userID, MyListKey, order, originalJSONOrder)
}

func (p *Plugin) storeIn(userID string, order []string, originalJSONOrder []byte) (bool, error) {
	return p.storeList(userID, InListKey, order, originalJSONOrder)
}

func (p *Plugin) storeOut(userID string, order []string, originalJSONOrder []byte) (bool, error) {
	return p.storeList(userID, OutListKey, order, originalJSONOrder)
}

func (p *Plugin) storeList(userID string, listID string, order []string, originalJSONOrder []byte) (bool, error) {
	newJSONOrder, jsonErr := json.Marshal(order)
	if jsonErr != nil {
		return false, jsonErr
	}

	ok, appErr := p.API.KVCompareAndSet(getListKey(userID, listID), originalJSONOrder, newJSONOrder)
	if appErr != nil {
		return false, errors.New(appErr.Error())
	}

	return ok, nil
}

func (p *Plugin) getOrderForUser(userID string) ([]string, []byte, error) {
	return p.getListForUser(userID, MyListKey)
}

func (p *Plugin) getInForUser(userID string) ([]string, []byte, error) {
	return p.getListForUser(userID, InListKey)
}

func (p *Plugin) getOutForUser(userID string) ([]string, []byte, error) {
	return p.getListForUser(userID, OutListKey)
}

func (p *Plugin) getListForUser(userID string, listID string) ([]string, []byte, error) {
	originalJSONOrder, err := p.API.KVGet(getListKey(userID, listID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONOrder == nil {
		return []string{}, originalJSONOrder, nil
	}

	var order []string
	jsonErr := json.Unmarshal(originalJSONOrder, &order)
	if jsonErr != nil {
		return nil, nil, jsonErr
	}

	return order, originalJSONOrder, nil
}

func getOrderKey(userID string) string {
	return getListKey(userID, MyListKey)
}

func getInKey(userID string) string {
	return getListKey(userID, InListKey)
}

func getOutKey(userID string) string {
	return getListKey(userID, OutListKey)
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

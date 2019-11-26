package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-server/model"
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
	order, _, err := p.getOrderForUser(userID)
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
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getOrderForUser(userID)
		if err != nil {
			return err
		}

		for _, id := range order {
			if id == itemID {
				return errors.New("item id already exists in order")
			}
		}

		order = append(order, itemID)

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

	return errors.New("unable to store installation")
}

func (p *Plugin) removeFromOrderForUser(userID string, itemID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getOrderForUser(userID)
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
	newJSONOrder, jsonErr := json.Marshal(order)
	if jsonErr != nil {
		return false, jsonErr
	}

	ok, appErr := p.API.KVCompareAndSet(getOrderKey(userID), originalJSONOrder, newJSONOrder)
	if appErr != nil {
		return false, errors.New(appErr.Error())
	}

	return ok, nil
}

func (p *Plugin) getOrderForUser(userID string) ([]string, []byte, error) {
	originalJSONOrder, err := p.API.KVGet(getOrderKey(userID))
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
	return fmt.Sprintf("%s_%s", StoreOrderKey, userID)
}

func getItemKey(itemID string) string {
	return fmt.Sprintf("%s_%s", StoreItemKey, itemID)
}

func getReminderKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreReminderKey, userID)
}

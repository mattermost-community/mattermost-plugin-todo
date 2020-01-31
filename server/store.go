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
)

// OrderElement denotes every element in any of the lists. Contains the item that refers to,
// and may contain foreign ids of item and user, denoting the user this element is related to
// and the item on that user system.
type OrderElement struct {
	ItemID        string `json:"item_id"`
	ForeignItemID string `json:"foreign_item_id"`
	ForeignUserID string `json:"foreign_user_id"`
}

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

func (p *Plugin) getItemListForUser(userID string, listID string) ([]*Item, error) {
	order, _, err := p.getListForUser(userID, listID)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for _, oe := range order {
		item, _, err := p.getItem(oe.ItemID)
		if err != nil {
			return nil, err
		}
		if item != nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func (p *Plugin) addToListForUser(userID string, itemID string, listID string, foreignItemID string, foreignUserID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getListForUser(userID, listID)
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

func (p *Plugin) removeFromListForUser(userID string, itemID string, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getListForUser(userID, listID)
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

func (p *Plugin) popFromMyListForUser(userID string) error {
	for i := 0; i < StoreRetries; i++ {
		order, originalJSONOrder, err := p.getListForUser(userID, MyListKey)
		if err != nil {
			return err
		}

		if len(order) == 0 {
			return nil
		}

		order = order[1:]

		ok, err := p.storeList(userID, MyListKey, order, originalJSONOrder)
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

func (p *Plugin) storeList(userID string, listID string, order []*OrderElement, originalJSONOrder []byte) (bool, error) {
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

func (p *Plugin) getListForUser(userID string, listID string) ([]*OrderElement, []byte, error) {
	originalJSONOrder, err := p.API.KVGet(getListKey(userID, listID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONOrder == nil {
		return []*OrderElement{}, originalJSONOrder, nil
	}

	var order []*OrderElement
	jsonErr := json.Unmarshal(originalJSONOrder, &order)
	if jsonErr != nil {
		return p.legacyOrderElement(userID, listID)
	}

	return order, originalJSONOrder, nil
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

func (p *Plugin) legacyOrderElement(userID string, listID string) ([]*OrderElement, []byte, error) {
	originalJSONOrder, err := p.API.KVGet(getListKey(userID, listID))
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

func (p *Plugin) getOrderForItem(userID string, itemID string, listID string) (*OrderElement, int, error) {
	originalJSONOrder, err := p.API.KVGet(getListKey(userID, listID))
	if err != nil {
		return nil, 0, err
	}

	if originalJSONOrder == nil {
		return nil, 0, nil
	}

	var order []*OrderElement
	jsonErr := json.Unmarshal(originalJSONOrder, &order)
	if jsonErr != nil {
		order, _, jsonErr = p.legacyOrderElement(userID, listID)
		if order == nil {
			return nil, 0, jsonErr
		}
	}

	for i, oe := range order {
		if oe.ItemID == itemID {
			return oe, i, nil
		}
	}
	return nil, 0, nil
}

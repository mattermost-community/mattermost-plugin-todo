package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	// WSEventRefresh is the WebSocket event for refreshing the to do list
	WSEventRefresh = "refresh"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	BotUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		return err
	}

	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "todo",
		DisplayName: "To Do Bot",
		Description: "Created by the To Do plugin.",
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure todo bot")
	}
	p.BotUserID = botID

	return p.API.RegisterCommand(getCommand())
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/add":
		p.handleAdd(w, r)
	case "/list":
		p.handleList(w, r)
	case "/remove":
		p.handleRemove(w, r)
	case "/complete":
		p.handleComplete(w, r)
	case "/enqueue":
		p.handleEnqueue(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleAdd(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var item *Item
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item.ID = model.NewId()
	item.CreateAt = model.GetMillis()

	p.storeItem(item)
	myList := p.getMyListForUser(userID)

	err = myList.add(item.ID, "", "")
	if err != nil {
		p.deleteItem(item.ID)
		p.API.LogError("Unable to add item err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) handleList(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	list := p.getMyListForUser(userID)
	listInput := r.URL.Query().Get("list")
	switch listInput {
	case "out":
		list = p.getOutListForUser(userID)
	case "in":
		list = p.getInListForUser(userID)
	}

	items, err := list.getItems()
	if err != nil {
		p.API.LogError("Unable to get items for user err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(items) > 0 && r.URL.Query().Get("reminder") == "true" {
		var lastReminderAt int64
		lastReminderAt, err = p.getLastReminderTimeForUser(userID)
		if err != nil {
			p.API.LogError("Unable to send reminder err=" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var timezone *time.Location
		offset, _ := strconv.Atoi(r.Header.Get("X-Timezone-Offset"))
		timezone = time.FixedZone("local", -60*offset)

		// Post reminder message if it's the next day and been more than an hour since the last post
		now := model.GetMillis()
		nt := time.Unix(now/1000, 0).In(timezone)
		lt := time.Unix(lastReminderAt/1000, 0).In(timezone)
		if nt.Sub(lt).Hours() >= 1 && (nt.Day() != lt.Day() || nt.Month() != lt.Month() || nt.Year() != lt.Year()) {
			p.PostBotDM(userID, "Daily Reminder:\n\n"+itemsListToString(items))
			p.storeLastReminderTimeForUser(userID)
		}
	}

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		p.API.LogError("Unable marhsal items list to json err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(itemsJSON)
}

type enqueueAPIRequest struct {
	ID string
}

func (p *Plugin) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var enqueueRequest *enqueueAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&enqueueRequest)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, _, err := p.getItem(enqueueRequest.ID)
	if err != nil {
		p.API.LogError("Unable to enqueue item err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	myList := p.getMyListForUser(userID)
	inList := p.getMyListForUser(userID)

	oe, _, err := inList.getOrderForItem(enqueueRequest.ID)
	if err != nil {
		p.API.LogError("Unable to enqueue item err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if oe == nil {
		p.API.LogError("Unable to enqueue item")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = myList.add(oe.ItemID, oe.ForeignItemID, oe.ForeignUserID)
	if err != nil {
		p.API.LogError("Unable to enqueue item err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = inList.remove(enqueueRequest.ID)
	if err != nil {
		myList.remove(oe.ItemID)
		p.API.LogError("Unable to enqueue item err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("%s enqueued a Todo you sent: %s", userID, item.Message)
	p.PostBotDM(oe.ForeignUserID, message)
}

type completeAPIRequest struct {
	ID string
}

func (p *Plugin) handleComplete(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var completeRequest *completeAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&completeRequest)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	itemList := p.getMyListForUser(userID)

	itemList, oe, _, _ := p.getUserListForItem(userID, completeRequest.ID)
	if itemList == nil {
		p.API.LogError("Unable to get item to complete")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := itemList.remove(oe.ItemID); err != nil {
		p.API.LogError("Unable to complete the item err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.deleteItem(oe.ItemID)
	p.handleExternalComplete(oe.ForeignUserID, oe.ForeignItemID, userID)
}

func (p *Plugin) handleExternalComplete(foreignUserID string, foreignItemID string, userID string) {
	if foreignUserID == "" {
		return
	}

	outList := p.getOutListForUser(foreignUserID)
	outList.remove(foreignItemID)
	item, _, err := p.getItem(foreignItemID)
	if err != nil {
		return
	}

	p.deleteItem(foreignItemID)
	message := fmt.Sprintf("%s completed a Todo you sent: %s", userID, item.Message)
	p.PostBotDM(userID, message)
}

type removeAPIRequest struct {
	ID string
}

func (p *Plugin) handleRemove(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var removeRequest *removeAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&removeRequest)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	itemList, oe, _, _ := p.getUserListForItem(userID, removeRequest.ID)
	if itemList == nil {
		p.API.LogError("Unable to get item to remove")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := itemList.remove(oe.ItemID); err != nil {
		p.API.LogError("Unable to remove the item err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.deleteItem(oe.ItemID)
	p.handleExternalRemove(oe.ForeignUserID, oe.ForeignItemID, userID)
}

func (p *Plugin) handleExternalRemove(foreignUserID string, foreignItemID string, userID string) {
	if foreignUserID == "" {
		return
	}

	itemList, _, _, listKey := p.getUserListForItem(foreignUserID, foreignItemID)

	itemList.remove(foreignItemID)
	item, _, err := p.getItem(foreignItemID)
	if err != nil {
		return
	}

	p.deleteItem(foreignItemID)
	message := fmt.Sprintf("%s removed a Todo you received: %s", userID, item.Message)
	if listKey == OutListKey {
		message = fmt.Sprintf("%s declined a Todo you sent: %s", userID, item.Message)
	}
	p.PostBotDM(userID, message)
}

func (p *Plugin) sendRefreshEvent(userID string) {
	p.API.PublishWebSocketEvent(
		WSEventRefresh,
		nil,
		&model.WebsocketBroadcast{UserId: userID},
	)
}

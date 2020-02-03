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

// ListManager representes the logic on the lists
type ListManager interface {
	Add(userID string, message string) error
	Send(senderID string, receiverID string, message string) error
	Get(userID string, listID string) ([]*ExtendedItem, error)
	Complete(userID string, itemID string) (todoMessage string, foreignUserID string, err error)
	Enqueue(userID string, itemID string) (todoMessage string, foreignUserID string, err error)
	Remove(userID string, itemID string) (todoMessage string, foreignUserID string, isSender bool, err error)
	Pop(userID string) (todoMessage string, sender string, err error)
}

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	BotUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	listManager ListManager
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

	p.listManager = NewListManager(NewListStore(p.API), p.getUserName)

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

	if err = p.listManager.Add(userID, item.Message); err != nil {
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

	listInput := r.URL.Query().Get("list")
	listID := MyListKey
	switch listInput {
	case "out":
		listID = OutListKey
	case "in":
		listID = InListKey
	}

	items, err := p.listManager.Get(userID, listID)
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
			p.saveLastReminderTimeForUser(userID)
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
	if err := decoder.Decode(&enqueueRequest); err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todoMessage, sender, err := p.listManager.Enqueue(userID, enqueueRequest.ID)

	if err != nil {
		p.API.LogError("Unable to enqueue item err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userName := p.getUserName(userID)

	message := fmt.Sprintf("%s enqueued a Todo you sent: %s", userName, todoMessage)
	p.PostBotDM(sender, message)
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
	if err := decoder.Decode(&completeRequest); err != nil {
		p.API.LogError("unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todoMessage, sender, err := p.listManager.Complete(userID, completeRequest.ID)
	if err != nil {
		p.API.LogError("unable to complete item, err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if sender == "" {
		return
	}

	userName := p.getUserName(sender)

	message := fmt.Sprintf("%s completed a Todo you sent: %s", userName, todoMessage)
	p.sendRefreshEvent(sender)
	p.PostBotDM(sender, message)
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

	todoMessage, foreignUser, isSender, err := p.listManager.Remove(userID, removeRequest.ID)
	if err != nil {
		p.API.LogError("unable to complete item, err=" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if foreignUser == "" {
		return
	}

	userName := p.getUserName(userID)

	message := fmt.Sprintf("%s removed a Todo you received: %s", userName, todoMessage)
	if isSender {
		message = fmt.Sprintf("%s declined a Todo you sent: %s", userName, todoMessage)
	}

	p.sendRefreshEvent(foreignUser)
	p.PostBotDM(foreignUser, message)
}

func (p *Plugin) sendRefreshEvent(userID string) {
	p.API.PublishWebSocketEvent(
		WSEventRefresh,
		nil,
		&model.WebsocketBroadcast{UserId: userID},
	)
}

func (p *Plugin) getUserName(userID string) string {
	userName := "Someone"
	if user, err := p.API.GetUser(userID); err == nil {
		userName = user.Username
	}
	return userName
}

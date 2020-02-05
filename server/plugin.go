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
	// Add adds a todo to userID's myList with the message
	Add(userID, message string) error
	// Send sends the todo with the message from senderID to receiverID and returns the receiver's itemID
	Send(senderID, receiverID, message string) (string, error)
	// Get gets the todos on listID for userID
	Get(userID, listID string) ([]*ExtendedItem, error)
	// Complete completes the todo itemID for userID, and returns the message and the foreignUserID if any
	Complete(userID, itemID string) (todoMessage string, foreignUserID string, err error)
	// Enqueue moves one the todo itemID of userID from inbox to myList, and returns the message and the foreignUserID if any
	Enqueue(userID, itemID string) (todoMessage string, foreignUserID string, err error)
	// Remove removes the todo itemID for userID and returns the message, the foreignUserID if any, and whether the user sent the todo to someone else
	Remove(userID, itemID string) (todoMessage string, foreignUserID string, isSender bool, err error)
	// Removes the first element of myList for userID and returns the message and the sender of that todo if any
	Pop(userID string) (todoMessage string, sender string, err error)
	// GetUserName returns the readable username from userID
	GetUserName(userID string) string
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

	p.listManager = NewListManager(p.API)

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

type addAPIRequest struct {
	Message string
	SendTo  string
}

func (p *Plugin) handleAdd(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var addRequest *addAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&addRequest)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	if addRequest.SendTo == "" {
		if err = p.listManager.Add(userID, addRequest.Message); err != nil {
			p.API.LogError("Unable to add item err=" + err.Error())
			p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to add item", err)
		}
		return
	}

	itemID, err := p.listManager.Send(userID, addRequest.SendTo, addRequest.Message)

	if err != nil {
		p.API.LogError("Unable to send item err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to send item", err)
		return
	}

	senderName := p.listManager.GetUserName(userID)

	receiverMessage := fmt.Sprintf("You have received a new Todo from @%s", senderName)

	p.PostBotCustomDM(addRequest.SendTo, receiverMessage, addRequest.Message, itemID)
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
		p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to get items for user", err)
		return
	}

	if len(items) > 0 && r.URL.Query().Get("reminder") == "true" {
		var lastReminderAt int64
		lastReminderAt, err = p.getLastReminderTimeForUser(userID)
		if err != nil {
			p.API.LogError("Unable to send reminder err=" + err.Error())
			p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to send reminder", err)
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
		p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable marhsal items list to json", err)
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
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	todoMessage, sender, err := p.listManager.Enqueue(userID, enqueueRequest.ID)

	if err != nil {
		p.API.LogError("Unable to enqueue item err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to enqueue item", err)
		return
	}

	userName := p.listManager.GetUserName(userID)

	message := fmt.Sprintf("@%s enqueued a Todo you sent: %s", userName, todoMessage)
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
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	todoMessage, sender, err := p.listManager.Complete(userID, completeRequest.ID)
	if err != nil {
		p.API.LogError("Unable to complete item err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to complete item", err)
		return
	}

	if sender == "" {
		return
	}

	userName := p.listManager.GetUserName(sender)

	message := fmt.Sprintf("@%s completed a Todo you sent: %s", userName, todoMessage)
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
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	todoMessage, foreignUser, isSender, err := p.listManager.Remove(userID, removeRequest.ID)
	if err != nil {
		p.API.LogError("Unable to remove item, err=" + err.Error())
		p.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to remove item", err)
		return
	}

	if foreignUser == "" {
		return
	}

	userName := p.listManager.GetUserName(userID)

	message := fmt.Sprintf("@%s removed a Todo you received: %s", userName, todoMessage)
	if isSender {
		message = fmt.Sprintf("@%s declined a Todo you sent: %s", userName, todoMessage)
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

func (p *Plugin) handleErrorWithCode(w http.ResponseWriter, code int, errTitle string, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(struct {
		Error   string `json:"error"`
		Details string `json:"details"`
	}{
		Error:   errTitle,
		Details: err.Error(),
	})
	_, _ = w.Write(b)
}

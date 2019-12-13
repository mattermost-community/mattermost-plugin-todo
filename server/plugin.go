package main

import (
	"encoding/json"
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

	err = p.storeItemForUser(userID, item)
	if err != nil {
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

	items, err := p.getItemsForUser(userID)
	if err != nil {
		p.API.LogError("Unable to get items for user err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(items) > 0 && r.URL.Query().Get("reminder") == "true" {
		lastReminderAt, err := p.getLastReminderTimeForUser(userID)
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

	err = p.removeFromOrderForUser(userID, removeRequest.ID)
	if err != nil {
		p.API.LogError("Unable to remove item err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) sendRefreshEvent(userID string) {
	p.API.PublishWebSocketEvent(
		WSEventRefresh,
		nil,
		&model.WebsocketBroadcast{UserId: userID},
	)
}

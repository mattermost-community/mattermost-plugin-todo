package plugin

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-plugin-todo/server/api"
	"github.com/mattermost/mattermost-plugin-todo/server/command"
	"github.com/mattermost/mattermost-plugin-todo/server/todo"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/pkg/errors"
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

	todo *todo.ListManager
	api  *api.Service
}

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		return err
	}

	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "todo",
		DisplayName: "Todo Bot",
		Description: "Created by the Todo plugin.",
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure todo bot")
	}
	p.BotUserID = botID
	todo := todo.NewListManager(p.API)
	p.todo = todo

	router := &mux.Router{}
	p.api = api.NewService(router, p.API, todo, p, p)

	return p.API.RegisterCommand(command.GetCommand())
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	com := command.NewCommand(args, p.API, p, p.todo)
	responses := com.Handle()
	for _, response := range responses {
		if response == nil {
			continue
		}
		var err error
		switch response.Type {
		case command.ResponseTypeEphemeral:
			p.PostCommandResponse(response.UserID, response.ChannelID, response.Message)
		case command.ResponseTypeBotCustomDM:
			err = p.PostBotCustomDM(response.UserID, response.Message, response.Todo, response.IssueID)
		case command.ResponseTypeBotDM:
			err = p.PostBotDM(response.UserID, response.Message)
		case command.ResponseTypeBotReplyPost:
			err = p.ReplyPostBot(response.UserID, response.Message, response.Todo)
		}
		if err != nil {
			p.API.LogError(err.Error())
		}
	}
	return &model.CommandResponse{}, nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.api.ServeHTTP(w, r)
}

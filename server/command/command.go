package command

import (
	"encoding/json"

	"github.com/mattermost/mattermost-plugin-todo/server/todo"
	"github.com/mattermost/mattermost-server/v5/model"
)

type Command struct {
	Args      *model.CommandArgs
	Log       logger
	pluginAPI pluginAPI
	todoAPI   todoAPI
}

type ResponseType string

// ResponseType values
const (
	ResponseTypeEphemeral    ResponseType = "ephemeral"
	ResponseTypePost         ResponseType = "post"
	ResponseTypeBotCustomDM  ResponseType = "bot_custom_dm"
	ResponseTypeBotReplyPost ResponseType = "bot_reply_post"
	ResponseTypeBotDM        ResponseType = "bot_dm"
	ResponseTypePayload      ResponseType = "payload"
)

type Response struct {
	Type      ResponseType
	Message   string
	Todo      string
	UserID    string
	ChannelID string
	IssueID   string
	PostID    string
}

type logger interface {
	LogError(msg string, keyValuePairs ...interface{})
}

type pluginAPI interface {
	GetUserByUsername(name string) (*model.User, error)
	SendRefreshEvent(userID string)
}

type todoAPI interface {
	SendIssue(senderID, receiverID, message, postID string) (string, error)
	GetUserName(userID string) string
	AddIssue(userID, message, postID string) error
	GetIssueList(userID, listID string) (todo.ExtendedIssues, error)
	PopIssue(userID string) (issue *todo.Issue, foreignID string, err error)
}

func NewCommand(args *model.CommandArgs, log logger, api pluginAPI, todo todoAPI) *Command {
	return &Command{
		Args:      args,
		Log:       log,
		pluginAPI: api,
		todoAPI:   todo,
	}
}

func GetCommand() *model.Command {
	return &model.Command{
		Trigger:          "todo",
		DisplayName:      "Todo Bot",
		Description:      "Interact with your Todo list.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: add, list, pop, send, help",
		AutoCompleteHint: "[command]",
	}
}

func getHelp() string {
	return `Available Commands:

add [message]
	Adds a Todo.

	example: /todo add Don't forget to be awesome

list
	Lists your Todo issues.

list [listName]
	List your issues in certain list

	example: /todo list in
	example: /todo list out
	example (same as /todo list): /todo list my

pop
	Removes the Todo issue at the top of the list.

send [user] [message]
	Sends some user a Todo

	example: /todo send @awesomePerson Don't forget to be awesome

help
	Display usage.
`
}

func (r *Response) ToJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

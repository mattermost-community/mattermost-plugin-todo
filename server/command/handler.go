package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-todo/server/todo"
	"github.com/mattermost/mattermost-server/v5/model"
)

// Handle method handles command and returns list of responses.
// Each command can have multiple responses. For example command can
// post ephemeral message as well DM user.
// Caller of the Handle what should be done with the Responses, one
// might want to output the response to the user or send them down the chain
// for example to another plugin.
func (c *Command) Handle() []*Response {
	stringArgs := strings.Split(strings.TrimSpace(c.Args.Command), " ")
	lengthOfArgs := len(stringArgs)
	restOfArgs := []string{}

	var handler func([]string, *model.CommandArgs) []*Response
	if lengthOfArgs == 1 {
		handler = c.runListCommand
	} else {
		command := stringArgs[1]
		if lengthOfArgs > 2 {
			restOfArgs = stringArgs[2:]
		}
		switch command {
		case "add":
			handler = c.runAddCommand
		case "list":
			handler = c.runListCommand
		case "pop":
			handler = c.runPopCommand
		case "send":
			handler = c.runSendCommand
		default:
			return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: getHelp()}}
		}
	}
	return handler(restOfArgs, c.Args)
}

func (c *Command) runSendCommand(args []string, extra *model.CommandArgs) []*Response {
	if len(args) < 2 {
		message := "You must specify a user and a message.\n" + getHelp()
		return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: message}}
	}

	userName := args[0]
	if args[0][0] == '@' {
		userName = args[0][1:]
	}
	receiver, appErr := c.pluginAPI.GetUserByUsername(userName)
	if appErr != nil {
		message := "Please, provide a valid user.\n" + getHelp()
		return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: message}}
	}

	if receiver.Id == extra.UserId {
		return c.runAddCommand(args[1:], extra)
	}

	message := strings.Join(args[1:], " ")

	receiverIssueID, err := c.todoAPI.SendIssue(extra.UserId, receiver.Id, message, "")
	if err != nil {
		return c.unknownError(err)
	}

	c.pluginAPI.SendRefreshEvent(extra.UserId)
	c.pluginAPI.SendRefreshEvent(receiver.Id)

	responseMessage := fmt.Sprintf("Todo sent to @%s.", userName)

	senderName := c.todoAPI.GetUserName(extra.UserId)

	receiverMessage := fmt.Sprintf("You have received a new Todo from @%s", senderName)

	botResponse := &Response{
		Type:    ResponseTypeBotCustomDM,
		Message: receiverMessage,
		Todo:    message,
		UserID:  receiver.Id,
		IssueID: receiverIssueID,
	}
	ephemeralResponse := &Response{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: responseMessage}

	return []*Response{botResponse, ephemeralResponse}
}

func (c *Command) unknownError(err error) []*Response {
	c.Log.LogError(err.Error())
	message := "An unknown error occurred. Please talk to your system administrator for help."
	return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: message}}
}

func (c *Command) runAddCommand(args []string, extra *model.CommandArgs) []*Response {
	message := strings.Join(args, " ")

	if message == "" {
		message = "Please add a task."
		return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: message}}
	}

	if err := c.todoAPI.AddIssue(extra.UserId, message, ""); err != nil {
		return c.unknownError(err)
	}

	c.pluginAPI.SendRefreshEvent(extra.UserId)

	responseMessage := "Added Todo."

	issues, err := c.todoAPI.GetIssueList(extra.UserId, todo.MyListKey)
	if err != nil {
		c.Log.LogError(err.Error())
		return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: responseMessage}}
	}

	responseMessage += " Todo List:\n\n"
	responseMessage += issues.ToString()
	return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: responseMessage}}
}

func (c *Command) runListCommand(args []string, extra *model.CommandArgs) []*Response {
	listID := todo.MyListKey
	responseMessage := "Todo List:\n\n"

	if len(args) > 0 {
		switch args[0] {
		case "my":
		case "in":
			listID = todo.InListKey
			responseMessage = "Received Todo list:\n\n"
		case "out":
			listID = todo.OutListKey
			responseMessage = "Sent Todo list:\n\n"
		default:
			return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: getHelp()}}
		}
	}

	issues, err := c.todoAPI.GetIssueList(extra.UserId, listID)
	if err != nil {
		return c.unknownError(err)
	}
	c.pluginAPI.SendRefreshEvent(extra.UserId)

	responseMessage += issues.ToString()
	return []*Response{{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: responseMessage}}
}

func (c *Command) runPopCommand(args []string, extra *model.CommandArgs) []*Response {
	issue, foreignID, err := c.todoAPI.PopIssue(extra.UserId)
	if err != nil {
		return c.unknownError(err)
	}

	userName := c.todoAPI.GetUserName(extra.UserId)

	var botDMResponse *Response
	if foreignID != "" {
		message := fmt.Sprintf("@%s popped a Todo you sent: %s", userName, issue.Message)
		c.pluginAPI.SendRefreshEvent(foreignID)
		botDMResponse = &Response{Type: ResponseTypeBotDM, UserID: foreignID, Message: message}
	}

	c.pluginAPI.SendRefreshEvent(extra.UserId)

	responseMessage := "Removed top Todo."

	replyMessage := fmt.Sprintf("@%s popped a todo attached to this thread", userName)
	botReplyResponse := &Response{Type: ResponseTypeBotReplyPost, PostID: issue.PostID, Message: replyMessage, Todo: issue.Message}

	var ephemeralResponse *Response
	issues, err := c.todoAPI.GetIssueList(extra.UserId, todo.MyListKey)
	if err != nil {
		c.Log.LogError(err.Error())
		ephemeralResponse = &Response{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: responseMessage}
	} else {
		responseMessage += " Todo List:\n\n"
		responseMessage += issues.ToString()
		ephemeralResponse = &Response{Type: ResponseTypeEphemeral, UserID: c.Args.UserId, ChannelID: c.Args.ChannelId, Message: responseMessage}
	}

	return []*Response{botDMResponse, botReplyResponse, ephemeralResponse}
}

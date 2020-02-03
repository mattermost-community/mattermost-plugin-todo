package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getHelp() string {
	return `Available Commands:

add [message]
	Adds a to do.

	example: /todo add Don't forget to be awesome

list
	Lists your to do items.

list [listName]
	List your items in certain list

	example: /todo list in
	example: /todo list out
	example (same as /todo list): /todo list my

pop
	Removes the to do item at the top of the list.

send [user] [message]
	Sends some user a Todo

	example: /todo send @awesomePerson Don't forget to be awesome

help
	Display usage.
`
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "todo",
		DisplayName:      "To Do Bot",
		Description:      "Interact with your to do list.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: add, list, pop",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "todo",
		//IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.ID),
	}
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	stringArgs := strings.Split(args.Command, " ")

	if len(stringArgs) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	command := stringArgs[1]

	var handler func([]string, *model.CommandArgs) (*model.CommandResponse, bool, error)

	switch command {
	case "add":
		handler = p.runAddCommand
	case "list":
		handler = p.runListCommand
	case "pop":
		handler = p.runPopCommand
	case "send":
		handler = p.runSendCommand
	}

	if handler == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	resp, isUserError, err := handler(stringArgs[2:], args)

	if err != nil {
		if isUserError {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("__Error: %s__\n\nRun `/todo help` for usage instructions.", err.Error())), nil
		}
		p.API.LogError(err.Error())
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "An unknown error occurred. Please talk to your system administrator for help."), nil
	}

	return resp, nil
}

func (p *Plugin) runSendCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "You must specify a user and a message."), false, nil
	}

	receiver, err := p.API.GetUserByUsername(args[0][1:])
	if err != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please, provide a valid user."), false, nil
	}

	if receiver.Id == extra.UserId {
		return p.runAddCommand(args[1:], extra)
	}

	message := strings.Join(args[1:], " ")

	if err := p.listManager.Send(extra.UserId, receiver.Id, message); err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)
	p.sendRefreshEvent(receiver.Id)

	responseMessage := fmt.Sprintf("Todo sent to %s.", args[0])

	senderName := p.getUserName(extra.UserId)

	receiverMessage := fmt.Sprintf("You have received a new Todo from %s: %s", senderName, message)

	p.PostBotDM(receiver.Id, receiverMessage)
	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

func (p *Plugin) runAddCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	message := strings.Join(args, " ")

	if message == "" {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please add a task."), false, nil
	}

	if err := p.listManager.Add(extra.UserId, message); err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)

	responseMessage := "Added to do."

	items, err := p.listManager.Get(extra.UserId, MyListKey)
	if err != nil {
		p.API.LogError(err.Error())
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
	}

	responseMessage += "To Do List:\n\n"
	responseMessage += itemsListToString(items)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

func (p *Plugin) runListCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	listID := MyListKey
	responseMessage := "To Do List:\n\n"

	if len(args) > 0 {
		switch args[0] {
		case "my":
		case "in":
			listID = InListKey
			responseMessage = "Received To Do list:\n\n"
		case "out":
			listID = OutListKey
			responseMessage = "Sent To Do list:\n\n"
		default:
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), true, nil
		}
	}

	items, err := p.listManager.Get(extra.UserId, listID)
	if err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)

	responseMessage += itemsListToString(items)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

func (p *Plugin) runPopCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	todoMessage, sender, err := p.listManager.Pop(extra.UserId)

	if sender != "" {
		userName := p.getUserName(sender)

		message := fmt.Sprintf("%s popped a Todo you sent: %s", userName, todoMessage)
		p.sendRefreshEvent(sender)
		p.PostBotDM(sender, message)
	}

	if err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)

	responseMessage := "Removed top to do."

	items, err := p.listManager.Get(extra.UserId, MyListKey)
	if err != nil {
		p.API.LogError(err.Error())
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
	}

	responseMessage += "To Do List:\n\n"
	responseMessage += itemsListToString(items)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

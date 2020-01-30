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

	example: /todo list inbox
	example: /todo list sent
	example (same as /todo list): /todo list own

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
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "You cannot send Todos to yourself. Use `/todo add` for this."), false, nil
	}

	message := strings.Join(args[1:], "")

	item := &Item{
		ID:       model.NewId(),
		CreateAt: model.GetMillis(),
		CreateBy: extra.UserId,
		SendTo:   receiver.Id,
		Status:   StatusPending,
		Message:  message,
	}

	appErr := p.storeSentItemForUsers(extra.UserId, receiver.Id, item)
	if appErr != nil {
		return nil, false, appErr
	}

	p.sendRefreshEvent(extra.UserId)
	p.sendRefreshEvent(receiver.Id)

	responseMessage := fmt.Sprintf("Todo sent to %s.", args[0])

	receiverMessage := fmt.Sprintf("You have received a new Todo from %s: %s", args[0], message)

	p.PostBotDM(receiver.Id, receiverMessage)
	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

func (p *Plugin) runAddCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	message := strings.Join(args, " ")

	if message == "" {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please add a task."), false, nil
	}

	item := &Item{
		ID:       model.NewId(),
		CreateAt: model.GetMillis(),
		Message:  message,
	}

	err := p.storeItemForUser(extra.UserId, item)
	if err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)

	responseMessage := "Added to do."

	items, err := p.getItemsForUser(extra.UserId)
	if err != nil {
		p.API.LogError(err.Error())
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
	}

	responseMessage += "To Do List:\n\n"
	responseMessage += itemsListToString(items)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

func (p *Plugin) runListCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	listID := OwnListKey
	if len(args) > 0 {
		switch args[0] {
		case "inbox":
			listID = InboxListKey
		case "sent":
			listID = SentListKey
		}
	}

	items, err := p.getItemListForUser(extra.UserId, listID)
	if err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)

	responseMessage := "To Do List:\n\n"
	responseMessage += itemsListToString(items)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

func (p *Plugin) runPopCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	err := p.popFromOrderForUser(extra.UserId)
	if err != nil {
		return nil, false, err
	}

	p.sendRefreshEvent(extra.UserId)

	responseMessage := "Removed top to do."

	items, err := p.getItemsForUser(extra.UserId)
	if err != nil {
		p.API.LogError(err.Error())
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
	}

	responseMessage += "To Do List:\n\n"
	responseMessage += itemsListToString(items)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, responseMessage), false, nil
}

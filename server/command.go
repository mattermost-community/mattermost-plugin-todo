package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

func getHelp() string {
	return `Available Commands:

add [message]
	Adds a to do.

	example: /todo add Don't forget to be awesome

list
	Lists your to do items.

pop
	Removes the to do item at the top of the list.

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
		AutoCompleteDesc: "Available commands: add, list",
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
	items, err := p.getItemsForUser(extra.UserId)
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

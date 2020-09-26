package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	listHeaderMessage = " Todo List:\n\n"
	MyFlag            = "my"
	InFlag            = "in"
	OutFlag           = "out"
)

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

settings summary [on, off]
	Sets user preference on daily reminders

	example: /todo settings summary on

help
	Display usage.
`
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "todo",
		DisplayName:      "Todo Bot",
		Description:      "Interact with your Todo list.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: add, list, pop, send, help",
		AutoCompleteHint: "[command]",
		AutocompleteData: getAutocompleteData(),
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	spaceRegExp := regexp.MustCompile(`\s+`)
	trimmedArgs := spaceRegExp.ReplaceAllString(strings.TrimSpace(args.Command), " ")
	stringArgs := strings.Split(trimmedArgs, " ")
	lengthOfArgs := len(stringArgs)
	restOfArgs := []string{}

	var handler func([]string, *model.CommandArgs) (bool, error)
	if lengthOfArgs == 1 {
		handler = p.runListCommand
		p.trackCommand(args.UserId, "")
	} else {
		command := stringArgs[1]
		if lengthOfArgs > 2 {
			restOfArgs = stringArgs[2:]
		}
		switch command {
		case "add":
			handler = p.runAddCommand
		case "list":
			handler = p.runListCommand
		case "pop":
			handler = p.runPopCommand
		case "send":
			handler = p.runSendCommand
		case "settings":
			handler = p.runSettingsCommand
		default:
			if command == "help" {
				p.trackCommand(args.UserId, command)
			} else {
				p.trackCommand(args.UserId, "not found")
			}
			p.postCommandResponse(args, getHelp())
			return &model.CommandResponse{}, nil
		}
		p.trackCommand(args.UserId, command)
	}
	isUserError, err := handler(restOfArgs, args)
	if err != nil {
		if isUserError {
			p.postCommandResponse(args, fmt.Sprintf("__Error: %s.__\n\nRun `/todo help` for usage instructions.", err.Error()))
		} else {
			p.API.LogError(err.Error())
			p.postCommandResponse(args, "An unknown error occurred. Please talk to your system administrator for help.")
		}
	}

	return &model.CommandResponse{}, nil
}

func (p *Plugin) runSendCommand(args []string, extra *model.CommandArgs) (bool, error) {
	if len(args) < 2 {
		p.postCommandResponse(extra, "You must specify a user and a message.\n"+getHelp())
		return false, nil
	}

	userName := args[0]
	if args[0][0] == '@' {
		userName = args[0][1:]
	}
	receiver, appErr := p.API.GetUserByUsername(userName)
	if appErr != nil {
		p.postCommandResponse(extra, "Please, provide a valid user.\n"+getHelp())
		return false, nil
	}

	if receiver.Id == extra.UserId {
		return p.runAddCommand(args[1:], extra)
	}

	message := strings.Join(args[1:], " ")

	receiverIssueID, err := p.listManager.SendIssue(extra.UserId, receiver.Id, message, "")
	if err != nil {
		return false, err
	}

	p.trackSendIssue(extra.UserId, sourceCommand, false)

	go p.sendRefreshEvent(extra.UserId, []string{OutListKey})
	go p.sendRefreshEvent(receiver.Id, []string{InListKey})

	responseMessage := fmt.Sprintf("Todo sent to @%s.", userName)

	senderName := p.listManager.GetUserName(extra.UserId)

	receiverMessage := fmt.Sprintf("You have received a new Todo from @%s", senderName)

	p.PostBotCustomDM(receiver.Id, receiverMessage, message, receiverIssueID)
	p.postCommandResponse(extra, responseMessage)
	return false, nil
}

func (p *Plugin) runAddCommand(args []string, extra *model.CommandArgs) (bool, error) {
	message := strings.Join(args, " ")

	if message == "" {
		p.postCommandResponse(extra, "Please add a task.")
		return false, nil
	}

	newIssue, err := p.listManager.AddIssue(extra.UserId, message, "")
	if err != nil {
		return false, err
	}

	p.trackAddIssue(extra.UserId, sourceCommand, false)

	go p.sendRefreshEvent(extra.UserId, []string{MyListKey})

	responseMessage := "Added Todo."

	issues, err := p.listManager.GetIssueList(extra.UserId, MyListKey)
	if err != nil {
		p.API.LogError(err.Error())
		p.postCommandResponse(extra, responseMessage)
		return false, nil
	}

	// It's possible that database replication delay has resulted in the issue
	// list not containing the newly-added issue, so we check for that and
	// append the issue manually if necessary.
	var issueIncluded bool
	for _, issue := range issues {
		if newIssue.ID == issue.ID {
			issueIncluded = true
			break
		}
	}
	if !issueIncluded {
		issues = append(issues, &ExtendedIssue{
			Issue: *newIssue,
		})
	}

	responseMessage += listHeaderMessage
	responseMessage += issuesListToString(issues)
	p.postCommandResponse(extra, responseMessage)

	return false, nil
}

func (p *Plugin) runListCommand(args []string, extra *model.CommandArgs) (bool, error) {
	listID := MyListKey
	responseMessage := "Todo List:\n\n"

	if len(args) > 0 {
		switch args[0] {
		case MyFlag:
		case InFlag:
			listID = InListKey
			responseMessage = "Received Todo list:\n\n"
		case OutFlag:
			listID = OutListKey
			responseMessage = "Sent Todo list:\n\n"
		default:
			p.postCommandResponse(extra, getHelp())
			return true, nil
		}
	}

	issues, err := p.listManager.GetIssueList(extra.UserId, listID)
	if err != nil {
		return false, err
	}
	go p.sendRefreshEvent(extra.UserId, []string{MyListKey, OutListKey, InListKey})

	responseMessage += issuesListToString(issues)
	p.postCommandResponse(extra, responseMessage)

	return false, nil
}

func (p *Plugin) runPopCommand(args []string, extra *model.CommandArgs) (bool, error) {
	issue, foreignID, err := p.listManager.PopIssue(extra.UserId)
	if err != nil {
		if err.Error() == "cannot find issue" {
			p.postCommandResponse(extra, "There are no Todos to pop.")
			return false, nil
		}
		return false, err
	}

	userName := p.listManager.GetUserName(extra.UserId)

	if foreignID != "" {
		message := fmt.Sprintf("@%s popped a Todo you sent: %s", userName, issue.Message)
		go p.sendRefreshEvent(foreignID, []string{OutListKey})
		p.PostBotDM(foreignID, message)
	}

	go p.sendRefreshEvent(extra.UserId, []string{MyListKey})

	responseMessage := "Removed top Todo."

	replyMessage := fmt.Sprintf("@%s popped a todo attached to this thread", userName)
	p.postReplyIfNeeded(issue.PostID, replyMessage, issue.Message)

	issues, err := p.listManager.GetIssueList(extra.UserId, MyListKey)
	if err != nil {
		p.API.LogError(err.Error())
		p.postCommandResponse(extra, responseMessage)
		return false, nil
	}

	responseMessage += listHeaderMessage
	responseMessage += issuesListToString(issues)
	p.postCommandResponse(extra, responseMessage)

	return false, nil
}

func (p *Plugin) runSettingsCommand(args []string, extra *model.CommandArgs) (bool, error) {
	if len(args) < 1 {
		return true, errors.New("no setting selected")
	}

	if args[0] == "summary" {
		if len(args) < 2 {
			return true, errors.New("choose whether you want this setting `on` or `off`")
		}
		if len(args) > 2 {
			return true, errors.New("too many arguments")
		}
		var responseMessage string
		var err error

		switch args[1] {
		case "on":
			err = p.saveReminderPreference(extra.UserId, true)
			responseMessage = "You will start receiving daily summaries."
		case "off":
			err = p.saveReminderPreference(extra.UserId, false)
			responseMessage = "You will stop receiving daily summaries."
		default:
			responseMessage = "invalid input, allowed values for \"settings summary\" are `on` or `off`"
			return true, errors.New(responseMessage)
		}

		if err != nil {
			responseMessage = "error saving the reminder preference"
			p.API.LogDebug("runSettingsCommand: error saving the reminder preference", "error", err.Error())
			return false, errors.New(responseMessage)
		}

		p.postCommandResponse(extra, responseMessage)
	} else {
		return true, fmt.Errorf("setting `%s` not recognized", args[0])
	}

	return false, nil
}

func getAutocompleteData() *model.AutocompleteData {
	todo := model.NewAutocompleteData("todo", "[command]", "Available commands: list, add, pop, send, settings, help")

	add := model.NewAutocompleteData("add", "[message]", "Adds a Todo")
	add.AddTextArgument("E.g. be awesome", "[message]", "")
	todo.AddCommand(add)

	list := model.NewAutocompleteData("list", "[name]", "Lists your Todo issues")
	items := []model.AutocompleteListItem{{
		HelpText: "Received Todos",
		Hint:     "(optional)",
		Item:     "in",
	}, {
		HelpText: "Sent Todos",
		Hint:     "(optional)",
		Item:     "out",
	}}
	list.AddStaticListArgument("Lists your Todo issues", false, items)
	todo.AddCommand(list)

	pop := model.NewAutocompleteData("pop", "", "Removes the Todo issue at the top of the list")
	todo.AddCommand(pop)

	send := model.NewAutocompleteData("send", "[user] [todo]", "Sends a Todo to a specified user")
	send.AddTextArgument("Whom to send", "[@awesomePerson]", "")
	send.AddTextArgument("Todo message", "[message]", "")
	todo.AddCommand(send)

	settings := model.NewAutocompleteData("settings", "summary [on] [off]", "Sets the user settings")
	summary := model.NewAutocompleteData("summary", "[on] [off]", "Sets the summary settings")
	on := model.NewAutocompleteData("on", "", "sets the daily reminder to enable")
	off := model.NewAutocompleteData("off", "", "sets the daily reminder to disable")
	summary.AddCommand(on)
	summary.AddCommand(off)
	settings.AddCommand(summary)
	todo.AddCommand(settings)

	help := model.NewAutocompleteData("help", "", "Display usage")
	todo.AddCommand(help)
	return todo
}

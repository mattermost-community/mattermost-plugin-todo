package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

// PostBotDM posts a DM as the cloud bot user.
func (p *Plugin) PostBotDM(userID string, message string) error {
	channel, appError := p.API.GetDirectChannel(userID, p.BotUserID)
	if appError != nil {
		return appError
	}
	if channel == nil {
		return fmt.Errorf("could not get direct channel for bot and user_id=%s", userID)
	}

	_, appError = p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		ChannelId: channel.Id,
		Message:   message,
	})

	return appError
}

// PostBotCustomDM posts a DM as the cloud bot user using custom post with action buttons.
func (p *Plugin) PostBotCustomDM(userID string, message string, todo string, itemID string) error {
	channel, appError := p.API.GetDirectChannel(userID, p.BotUserID)
	if appError != nil {
		return appError
	}
	if channel == nil {
		return fmt.Errorf("could not get direct channel for bot and user_id=%s", userID)
	}

	_, appError = p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		ChannelId: channel.Id,
		Message:   message + ": " + todo,
		Type:      "custom_todo",
		Props: map[string]interface{}{
			"type":    "custom_todo",
			"message": message,
			"todo":    todo,
			"itemId":  itemID,
		},
	})

	return appError
}

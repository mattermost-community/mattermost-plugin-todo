package plugin

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/pkg/errors"
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
func (p *Plugin) PostBotCustomDM(userID string, message string, todo string, issueID string) error {
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
			"issueId": issueID,
		},
	})

	return appError
}

// ReplyPostBot post a message and a todo in the same thread as the post postID
func (p *Plugin) ReplyPostBot(postID, message, todo string) error {
	if postID == "" {
		return errors.New("Post ID not defined")
	}

	post, appErr := p.API.GetPost(postID)
	if appErr != nil {
		return appErr
	}
	rootID := post.Id
	if post.RootId != "" {
		rootID = post.RootId
	}

	quotedTodo := "\n> " + strings.Join(strings.Split(todo, "\n"), "\n> ")
	_, appErr = p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		ChannelId: post.ChannelId,
		Message:   message + quotedTodo,
		RootId:    rootID,
	})

	if appErr != nil {
		return appErr
	}

	return nil
}

func (p *Plugin) PostCommandResponse(userID, channelID, text string) {
	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: channelID,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(userID, post)
}

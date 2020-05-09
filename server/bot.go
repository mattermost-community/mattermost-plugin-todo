package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// PostBotDM posts a DM as the cloud bot user.
func (p *Plugin) PostBotDM(userID string, message string) {
	p.createBotPostDM(&model.Post{
		UserId:  p.BotUserID,
		Message: message,
	}, userID)
}

// PostBotCustomDM posts a DM as the cloud bot user using custom post with action buttons.
func (p *Plugin) PostBotCustomDM(userID string, message string, todo string, issueID string) {
	p.createBotPostDM(&model.Post{
		UserId:  p.BotUserID,
		Message: message + ": " + todo,
		Type:    "custom_todo",
		Props: map[string]interface{}{
			"type":    "custom_todo",
			"message": message,
			"todo":    todo,
			"issueId": issueID,
		},
	}, userID)
}

func (p *Plugin) createBotPostDM(post *model.Post, userID string) {
	channel, appError := p.API.GetDirectChannel(userID, p.BotUserID)

	if appError != nil {
		p.API.LogError("Unable to get direct channel for bot err=" + appError.Error())
		return
	}
	if channel == nil {
		p.API.LogError("Could not get direct channel for bot and user_id=%s", userID)
		return
	}

	post.ChannelId = channel.Id
	_, appError = p.API.CreatePost(post)

	if appError != nil {
		p.API.LogError("Unable to create bot post DM err=" + appError.Error())
	}
}

// ReplyPostBot post a message and a todo in the same thread as the post postID
func (p *Plugin) ReplyPostBot(postID, message, todo string) error {
	if postID == "" {
		return errors.New("post ID not defined")
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

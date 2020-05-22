// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-plugin-todo/server/todo"
	"github.com/mattermost/mattermost-server/v5/model"
)

type Service struct {
	*mux.Router
	log    logger
	todo   todoAPI
	plugin pluginAPI
	bot    botAPI
}

type logger interface {
	LogError(msg string, keyValuePairs ...interface{})
}

type pluginAPI interface {
	GetUserByUsername(name string) (*model.User, *model.AppError)
	SendRefreshEvent(userID string)
}

type botAPI interface {
	PostBotDM(userID string, message string) error
	ReplyPostBot(postID, message, todo string) error
	PostBotCustomDM(userID string, message string, todo string, issueID string) error
}

type todoAPI interface {
	SendIssue(senderID, receiverID, message, postID string) (string, error)
	GetUserName(userID string) string
	AddIssue(userID, message, postID string) error
	GetIssueList(userID, listID string) (todo.ExtendedIssues, error)
	AcceptIssue(userID, issueID string) (todoMessage string, foreignUserID string, outErr error)
	CompleteIssue(userID, issueID string) (issue *todo.Issue, foreignID string, err error)
	RemoveIssue(userID, issueID string) (outIssue *todo.Issue, foreignID string, isSender bool, outErr error)
	BumpIssue(userID, issueID string) (todoMessage string, receiver string, foreignIssueID string, outErr error)
	PopIssue(userID string) (issue *todo.Issue, foreignID string, err error)
	AddReminder(userID string) error
	GetReminder(userID string) (int64, error)
}

func NewService(router *mux.Router, log logger, todo todoAPI, plugin pluginAPI, bot botAPI) *Service {
	s := &Service{
		Router: router,
		log:    log,
		todo:   todo,
		plugin: plugin,
		bot:    bot,
	}
	router.HandleFunc("/add", s.handleAdd).Methods("POST")
	router.HandleFunc("/list", s.handleList).Methods("GET")
	router.HandleFunc("/remove", s.handleRemove).Methods("POST")
	router.HandleFunc("/complete", s.handleComplete).Methods("POST")
	router.HandleFunc("/accept", s.handleAccept).Methods("PST")
	router.HandleFunc("/bump", s.handleBump).Methods("POST")
	router.HandleFunc("/execute_command", s.executeCommand).Methods("POST")
	router.Handle("{anything:.*}", http.NotFoundHandler())
	return s
}

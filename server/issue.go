package main

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

// Issue represents a Todo issue
type Issue struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	CreateAt    int64  `json:"create_at"`
	PostID      string `json:"post_id"`
}

// ExtendedIssue extends the information on Issue to be used on the front-end
type ExtendedIssue struct {
	Issue
	ForeignUser     string `json:"user"`
	ForeignList     string `json:"list"`
	ForeignPosition int    `json:"position"`
}

func newIssue(message string, description, postID string) *Issue {
	return &Issue{
		ID:          model.NewId(),
		CreateAt:    model.GetMillis(),
		Message:     message,
		Description: description,
		PostID:      postID,
	}
}

func issuesListToString(issues []*ExtendedIssue) string {
	if len(issues) == 0 {
		return "Nothing to do!"
	}

	str := "\n\n"

	for _, issue := range issues {
		createAt := time.Unix(issue.CreateAt/1000, 0)
		str += fmt.Sprintf("* %s\n  * (%s)\n", issue.Message, createAt.Format("January 2, 2006 at 15:04"))
	}

	return str
}

package main

import (
	fbClient "github.com/mattermost/focalboard/server/client"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// ListStore represents the KVStore operations for lists
type FocalboardListStore interface {
	// Issue related function
	AddIssue(userID string, issue *Issue) error
	GetIssue(issueID string) (*Issue, error)
	RemoveIssue(issueID string) error
	GetAndRemoveIssue(issueID string) (*Issue, error)

	// Issue References related functions

	// AddReference creates a new IssueRef with the issueID, foreignUSerID and foreignIssueID, and stores it
	// on the listID for userID.
	AddReference(userID, issueID, listID, foreignUserID, foreignIssueID string) error
	// RemoveReference removes the IssueRef for issueID in listID for userID
	RemoveReference(userID, issueID, listID string) error
	// PopReference removes the first IssueRef in listID for userID and returns it
	PopReference(userID, listID string) (*IssueRef, error)
	// BumpReference moves the Issue reference for issueID in listID for userID to the beginning of the list
	BumpReference(userID, issueID, listID string) error

	// GetIssueReference gets the IssueRef and position of the issue issueID on user userID's list listID
	GetIssueReference(userID, issueID, listID string) (*IssueRef, int, error)
	// GetIssueListAndReference gets the issue list, IssueRef and position for user userID
	GetIssueListAndReference(userID, issueID string) (string, *IssueRef, int)

	// GetList returns the list of IssueRef in listID for userID
	GetList(userID, listID string) ([]*IssueRef, error)
}

type focalboardListManager struct {
	store FocalboardListStore
	api   plugin.API
}

// NewListManager creates a new listManager
func NewFocalboardListManager(api plugin.API, client *fbClient.Client) FocalboardListManager {
	return &focalboardListManager{
		store: NewFocalboardListStore(api, client),
		api:   api,
	}
}

func (l *focalboardListManager) AddIssue(userID, message, postID string) (*Issue, error) {
	issue := newIssue(message, postID)

	if err := l.store.AddIssue(userID, issue); err != nil {
		return nil, err
	}

	return issue, nil
}

func (l *focalboardListManager) SendIssue(senderID, receiverID, message, postID string) (string, error) {
	return "", nil
}

func (l *focalboardListManager) GetIssueList(userID, listID string) ([]*ExtendedIssue, error) {
	return []*ExtendedIssue{}, nil
}

func (l *focalboardListManager) CompleteIssue(userID, issueID string) (issue *Issue, foreignID string, listToUpdate string, err error) {
	return issue, "", "", nil
}

func (l *focalboardListManager) AcceptIssue(userID, issueID string) (todoMessage string, foreignUserID string, outErr error) {
	return "", "", nil
}

func (l *focalboardListManager) RemoveIssue(userID, issueID string) (outIssue *Issue, foreignID string, isSender bool, listToUpdate string, outErr error) {
	return &Issue{}, "", false, "", nil
}

func (l *focalboardListManager) PopIssue(userID string) (issue *Issue, foreignID string, err error) {
	return &Issue{}, "", nil
}

func (l *focalboardListManager) BumpIssue(userID, issueID string) (todoMessage string, receiver string, foreignIssueID string, outErr error) {
	return "", "", "", nil
}

func (l *focalboardListManager) GetUserName(userID string) string {
	user, err := l.api.GetUser(userID)
	if err != nil {
		return "Someone"
	}
	return user.Username
}

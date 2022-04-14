package main

import (
	fbClient "github.com/mattermost/focalboard/server/client"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"
)

// ListStore represents the KVStore operations for lists
type FocalboardListStore interface {
	// Issue related function
	AddIssue(userID string, issue *ExtendedIssue) error
	RemoveIssue(userID, issueID string) error
	GetIssue(userID, issueID string) (*ExtendedIssue, error)
	GetIssuesByListType(userID string) (map[string][]*ExtendedIssue, error)
	CompleteIssue(userID, issueID string) (*ExtendedIssue, string, error)

	// Issue References related functions

	// AddReference creates a new IssueRef with the issueID, foreignUSerID and foreignIssueID, and stores it
	// on the listID for userID.
	AddReference(issue *ExtendedIssue) error
	// RemoveReference removes the reference for issueID for userID
	RemoveReference(userID, issueID string) error

	// GetReferenceList returns the list of references for userID
	GetReferenceList(userID string) ([]*ExtendedIssue, []byte, error)
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

	err := l.store.AddIssue(userID, &ExtendedIssue{Issue: *issue})
	if err != nil {
		return nil, errors.Wrap(err, "unable to add issue")
	}

	return issue, nil
}

func (l *focalboardListManager) SendIssue(senderID, receiverID, message, postID string) (string, error) {
	issue := newIssue(message, postID)
	extendedIssue := &ExtendedIssue{Issue: *issue, ForeignUser: senderID}

	err := l.store.AddIssue(receiverID, extendedIssue)
	if err != nil {
		return "", errors.Wrap(err, "unable to create issue")
	}

	err = l.store.AddReference(extendedIssue)
	if err != nil {
		err = errors.Wrap(err, "unable to create issue")
		rollbackErr := l.store.RemoveIssue(receiverID, extendedIssue.ID)
		if rollbackErr != nil {
			err = errors.Wrap(err, "unable to rollback")
		}
		return "", err
	}

	return "", nil
}

// TODO fix inefficiency of getting two lists at once when we only need one
func (l *focalboardListManager) GetIssueList(userID, listID string) ([]*ExtendedIssue, error) {
	if listID == OutListKey {
		list, _, err := l.store.GetReferenceList(userID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get reference list")
		}
		return list, nil
	}

	lists, err := l.store.GetIssuesByListType(userID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get issue list")
	}

	return lists[listID], nil
}

func (l *focalboardListManager) CompleteIssue(userID, issueID string) (*Issue, string, string, error) {
	issue, listType, err := l.store.CompleteIssue(userID, issueID)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "unable to complete issue")
	}

	err = l.store.RemoveReference(issue.ForeignUser, issueID)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "unable to remove reference to complete issue")
	}

	return &issue.Issue, issue.ForeignUser, listType, nil
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

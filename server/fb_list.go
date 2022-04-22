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
	DeleteIssue(userID, issueID string) error
	GetIssue(userID, issueID string) (*ExtendedIssue, error)
	GetIssuesByListType(userID string) (map[string][]*ExtendedIssue, error)
	UpdateIssueStatus(userID, issueID, status string) (*ExtendedIssue, string, error)

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
	extendedIssue := &ExtendedIssue{Issue: *issue, ForeignUser: senderID, ForeignList: "fb"}

	err := l.store.AddIssue(receiverID, extendedIssue)
	if err != nil {
		return "", errors.Wrap(err, "unable to create issue")
	}

	err = l.store.AddReference(extendedIssue)
	if err != nil {
		err = errors.Wrap(err, "unable to create issue")
		rollbackErr := l.store.DeleteIssue(receiverID, extendedIssue.ID)
		if rollbackErr != nil {
			err = errors.Wrap(err, "unable to rollback")
		}
		return "", err
	}

	return "", nil
}

func (l *focalboardListManager) GetIssueList(userID, listID string) ([]*ExtendedIssue, error) {
	var list []*ExtendedIssue
	var err error
	if listID == OutListKey {
		list, _, err = l.store.GetReferenceList(userID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get reference list")
		}
	} else {
		// TODO fix inefficiency of getting two lists at once when we only need one
		lists, err := l.store.GetIssuesByListType(userID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get issue list")
		}
		list = lists[listID]
	}

	for _, issue := range list {
		issue.ForeignUser = l.GetUserName(issue.ForeignUser)
	}

	return list, nil
}

func (l *focalboardListManager) CompleteIssue(userID, issueID string) (*Issue, string, string, error) {
	issue, listType, err := l.store.UpdateIssueStatus(userID, issueID, StatusDone)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "unable to complete issue")
	}

	err = l.store.RemoveReference(issue.ForeignUser, issueID)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "unable to remove reference when completing issue")
	}

	return &issue.Issue, issue.ForeignUser, listType, nil
}

func (l *focalboardListManager) AcceptIssue(userID, issueID string) (todoMessage string, foreignUserID string, outErr error) {
	issue, _, err := l.store.UpdateIssueStatus(userID, issueID, StatusToDo)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to accept issue")
	}

	return issue.Message, issue.ForeignUser, nil
}

func (l *focalboardListManager) RemoveIssue(userID, issueID string) (*Issue, string, bool, string, error) {
	issue, listType, err := l.store.UpdateIssueStatus(userID, issueID, StatusWontDo)
	if err != nil {
		return nil, "", false, "", errors.Wrap(err, "unable to remove issue")
	}

	err = l.store.RemoveReference(issue.ForeignUser, issueID)
	if err != nil {
		return nil, "", false, "", errors.Wrap(err, "unable to remove reference when completing issue")
	}

	return &issue.Issue, issue.ForeignUser, false, listType, nil
}

func (l *focalboardListManager) PopIssue(userID string) (*Issue, string, error) {
	lists, err := l.store.GetIssuesByListType(userID)
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to get issue list for popping")
	}
	list := lists[MyListKey]

	if len(list) == 0 {
		return nil, "", nil
	}

	issue := list[0]

	_, _, err = l.store.UpdateIssueStatus(userID, issue.ID, StatusDone)
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to complete issue when popped")
	}

	err = l.store.RemoveReference(issue.ForeignUser, issue.ID)
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to remove reference when popping issue")
	}

	return &issue.Issue, issue.ForeignUser, nil
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

package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	// MyListKey is the key used to store the list of the owned todos
	MyListKey = ""
	// InListKey is the key used to store the list of received todos
	InListKey = "_in"
	// OutListKey is the key used to store the list of sent todos
	OutListKey = "_out"
)

// ListStore represents the KVStore operations for lists
type ListStore interface {
	// Issue related function
	AddIssue(issue *Issue) error
	GetIssue(issueID string) (*Issue, error)
	RemoveIssue(issueID string) (*Issue, error)

	// Issue References related functions

	// AddReference creates a new IssueRef with the issueID, foreignUSerID and foreignIssueID, and stores it
	// on the listID for userID.
	AddReference(userID, issueID, listID, foreignUserID, foreignIssueID string) error
	// RemoveReference removes the IssueRef for issueID in listID for userID
	RemoveReference(userID, issueID, listID string) error
	// PopReference removes the first IssueRef in listID for userID and returns it
	PopReference(userID, listID string) (*IssueRef, error)
	// BumpReference moves the Issue reference for issueID in listID for userID to the beggining of the list
	BumpReference(userID, issueID, listID string) error

	// GetIssueReference gets the IssueRef and position of the issue issueID on user userID's list listID
	GetIssueReference(userID, issueID, listID string) (*IssueRef, int, error)
	// GetIssueListAndReference gets the issue list, IssueRef and position for user userID
	GetIssueListAndReference(userID, issueID string) (string, *IssueRef, int)

	// GetList returns the list of IssueRef in listID for userID
	GetList(userID, listID string) ([]*IssueRef, error)
}

type listManager struct {
	store ListStore
	api   plugin.API
}

// NewListManager creates a new listManager
func NewListManager(api plugin.API) *listManager {
	return &listManager{
		store: NewListStore(api),
		api:   api,
	}
}

func (l *listManager) AddIssue(userID, message, postID string) error {
	issue := newIssue(message, postID)

	if err := l.store.AddIssue(issue); err != nil {
		return err
	}

	if err := l.store.AddReference(userID, issue.ID, MyListKey, "", ""); err != nil {
		if _, rollbackError := l.store.RemoveIssue(issue.ID); rollbackError != nil {
			l.api.LogError("cannot rollback issue after add error, Err=", err.Error())
		}
		return err
	}

	return nil
}

func (l *listManager) SendIssue(senderID, receiverID, message, postID string) (string, error) {
	senderIssue := newIssue(message, postID)
	if err := l.store.AddIssue(senderIssue); err != nil {
		return "", err
	}

	receiverIssue := newIssue(message, postID)
	if err := l.store.AddIssue(receiverIssue); err != nil {
		if _, rollbackError := l.store.RemoveIssue(senderIssue.ID); rollbackError != nil {
			l.api.LogError("cannot rollback sender issue after send error, Err=", err.Error())
		}
		return "", err
	}

	if err := l.store.AddReference(senderID, senderIssue.ID, OutListKey, receiverID, receiverIssue.ID); err != nil {
		if _, rollbackError := l.store.RemoveIssue(senderIssue.ID); rollbackError != nil {
			l.api.LogError("cannot rollback sender issue after send error, Err=", err.Error())
		}
		if _, rollbackError := l.store.RemoveIssue(receiverIssue.ID); rollbackError != nil {
			l.api.LogError("cannot rollback receiver issue after send error, Err=", err.Error())
		}
		return "", err
	}

	if err := l.store.AddReference(receiverID, receiverIssue.ID, InListKey, senderID, senderIssue.ID); err != nil {
		if _, rollbackError := l.store.RemoveIssue(senderIssue.ID); rollbackError != nil {
			l.api.LogError("cannot rollback sender issue after send error, Err=", err.Error())
		}
		if _, rollbackError := l.store.RemoveIssue(receiverIssue.ID); rollbackError != nil {
			l.api.LogError("cannot rollback receiver issue after send error ,Err=", err.Error())
		}
		if rollbackError := l.store.RemoveReference(senderID, senderIssue.ID, OutListKey); rollbackError != nil {
			l.api.LogError("cannot rollback sender list after send error, Err=", err.Error())
		}
		return "", err
	}

	return receiverIssue.ID, nil
}

func (l *listManager) GetIssueList(userID, listID string) ([]*ExtendedIssue, error) {
	irs, err := l.store.GetList(userID, listID)
	if err != nil {
		return nil, err
	}

	extendedIssues := []*ExtendedIssue{}
	for _, ir := range irs {
		issue, err := l.store.GetIssue(ir.IssueID)
		if err != nil {
			continue
		}

		extendedIssue := l.extendIssueInfo(issue, ir)
		extendedIssues = append(extendedIssues, extendedIssue)
	}

	return extendedIssues, nil
}

func (l *listManager) CompleteIssue(userID, issueID string) (*ExtendedIssue, error) {
	issueList, ir, _ := l.store.GetIssueListAndReference(userID, issueID)
	if ir == nil {
		return nil, fmt.Errorf("cannot find element")
	}

	if err := l.store.RemoveReference(userID, issueID, issueList); err != nil {
		return nil, err
	}

	issue, err := l.store.RemoveIssue(issueID)
	if err != nil {
		l.api.LogError("cannot remove issue, Err=", err.Error())
	}

	if ir.ForeignUserID == "" {
		if issue == nil {
			return &ExtendedIssue{}, nil
		}
		return l.extendIssueInfo(issue, ir), nil
	}

	err = l.store.RemoveReference(ir.ForeignUserID, ir.ForeignIssueID, OutListKey)
	if err != nil {
		l.api.LogError("cannot clean foreigner list after complete, Err=", err.Error())
	}

	issue, err = l.store.RemoveIssue(ir.ForeignIssueID)
	if err != nil {
		l.api.LogError("cannot clean foreigner issue after complete, Err=", err.Error())
	}

	return l.extendIssueInfo(issue, ir), nil
}

func (l *listManager) AcceptIssue(userID, issueID string) (todoMessage string, foreignUserID string, outErr error) {
	issue, err := l.store.GetIssue(issueID)
	if err != nil {
		return "", "", err
	}

	ir, _, err := l.store.GetIssueReference(userID, issueID, InListKey)
	if err != nil {
		return "", "", err
	}
	if ir == nil {
		return "", "", fmt.Errorf("element reference not found")
	}

	err = l.store.AddReference(userID, issueID, MyListKey, ir.ForeignUserID, ir.ForeignIssueID)
	if err != nil {
		return "", "", err
	}

	err = l.store.RemoveReference(userID, issueID, InListKey)
	if err != nil {
		if rollbackError := l.store.RemoveReference(userID, issueID, MyListKey); rollbackError != nil {
			l.api.LogError("cannot rollback accept operation, Err=", rollbackError.Error())
		}
		return "", "", err
	}

	return issue.Message, ir.ForeignUserID, nil
}

func (l *listManager) RemoveIssue(userID, issueID string) (outIssue *ExtendedIssue, isSender bool, outErr error) {
	issueList, ir, _ := l.store.GetIssueListAndReference(userID, issueID)
	if ir == nil {
		return nil, false, fmt.Errorf("cannot find element")
	}

	if err := l.store.RemoveReference(userID, issueID, issueList); err != nil {
		return nil, false, err
	}

	issue, err := l.store.RemoveIssue(issueID)
	if err != nil {
		l.api.LogError("cannot remove issue, Err=", err.Error())
	}

	if ir.ForeignUserID == "" {
		if issue == nil {
			return &ExtendedIssue{}, false, nil
		}
		return l.extendIssueInfo(issue, ir), false, nil
	}

	list, _, _ := l.store.GetIssueListAndReference(ir.ForeignUserID, ir.ForeignIssueID)

	err = l.store.RemoveReference(ir.ForeignUserID, ir.ForeignIssueID, list)
	if err != nil {
		l.api.LogError("cannot clean foreigner list after remove, Err=", err.Error())
	}

	issue, err = l.store.RemoveIssue(ir.ForeignIssueID)
	if err != nil {
		l.api.LogError("cannot clean foreigner issue after remove, Err=", err.Error())
	}

	return l.extendIssueInfo(issue, ir), list == OutListKey, nil
}

func (l *listManager) PopIssue(userID string) (*ExtendedIssue, error) {
	ir, err := l.store.PopReference(userID, MyListKey)
	if err != nil {
		return nil, err
	}

	if ir == nil {
		return &ExtendedIssue{}, nil
	}

	issue, err := l.store.RemoveIssue(ir.IssueID)
	if err != nil {
		l.api.LogError("cannot remove issue after pop, Err=", err.Error())
	}

	if ir.ForeignUserID == "" {
		if issue == nil {
			return &ExtendedIssue{}, nil
		}
		return l.extendIssueInfo(issue, ir), nil
	}

	err = l.store.RemoveReference(ir.ForeignUserID, ir.ForeignIssueID, OutListKey)
	if err != nil {
		l.api.LogError("cannot clean foreigner list after pop, Err=", err.Error())
	}
	issue, err = l.store.RemoveIssue(ir.ForeignIssueID)
	if err != nil {
		l.api.LogError("cannot clean foreigner issue after pop, Err=", err.Error())
	}

	return l.extendIssueInfo(issue, ir), nil
}

func (l *listManager) BumpIssue(userID, issueID string) (todoMessage string, receiver string, foreignIssueID string, outErr error) {
	ir, _, err := l.store.GetIssueReference(userID, issueID, OutListKey)
	if err != nil {
		return "", "", "", err
	}

	if ir == nil {
		return "", "", "", fmt.Errorf("cannot find sender issue")
	}

	err = l.store.BumpReference(ir.ForeignUserID, ir.ForeignIssueID, InListKey)
	if err != nil {
		return "", "", "", err
	}

	if ir == nil {
		return "", "", "", fmt.Errorf("cannot find receiver issue")
	}

	issue, err := l.store.GetIssue(ir.ForeignIssueID)
	if err != nil {
		l.api.LogError("cannot find foreigner issue after bump, Err=", err.Error())
		return "", "", "", nil
	}

	return issue.Message, ir.ForeignUserID, ir.ForeignIssueID, nil
}

func (l *listManager) GetUserName(userID string) string {
	user, err := l.api.GetUser(userID)
	if err != nil {
		return "Someone"
	}
	return user.Username
}

func (l *listManager) extendIssueInfo(issue *Issue, ir *IssueRef) *ExtendedIssue {
	if issue == nil || ir == nil {
		return nil
	}

	feIssue := &ExtendedIssue{
		Issue: *issue,
	}

	if ir == nil || ir.ForeignUserID == "" {
		return feIssue
	}

	list, _, n := l.store.GetIssueListAndReference(ir.ForeignUserID, ir.ForeignIssueID)

	var listName string
	switch list {
	case MyListKey:
		listName = ""
	case InListKey:
		listName = "in"
	case OutListKey:
		listName = "out"
	}

	userName := l.GetUserName(ir.ForeignUserID)

	feIssue.ForeignUser = userName
	feIssue.ForeignList = listName
	feIssue.ForeignPosition = n

	return feIssue
}

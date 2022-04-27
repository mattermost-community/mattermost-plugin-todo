package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	// StoreRetries is the number of retries to use when storing lists fails on a race
	StoreRetries = 3
	// StoreListKey is the key used to store lists in the plugin KV store. Still "order" for backwards compatibility.
	StoreListKey = "order"
	// StoreIssueKey is the key used to store issues in the plugin KV store. Still "item" for backwards compatibility.
	StoreIssueKey = "item"
	// StoreReminderKey is the key used to store the last time a user was reminded
	StoreReminderKey = "reminder"
	// StoreReminderEnabledKey is the key used to store the user preference of auto daily reminder
	StoreReminderEnabledKey = "reminder_enabled"

	// StoreAllowIncomingTaskRequestsKey is the key used to store user preference for wallowing any incoming todo requests
	StoreAllowIncomingTaskRequestsKey = "allow_incoming_task"
)

// IssueRef denotes every element in any of the lists. Contains the issue that refers to,
// and may contain foreign ids of issue and user, denoting the user this element is related to
// and the issue on that user system.
type IssueRef struct {
	IssueID        string `json:"issue_id"`
	ForeignIssueID string `json:"foreign_issue_id"`
	ForeignUserID  string `json:"foreign_user_id"`
}

func listKey(userID string, listID string) string {
	return fmt.Sprintf("%s_%s%s", StoreListKey, userID, listID)
}

func issueKey(issueID string) string {
	return fmt.Sprintf("%s_%s", StoreIssueKey, issueID)
}

func reminderKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreReminderKey, userID)
}

func reminderEnabledKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreReminderEnabledKey, userID)
}

func allowIncomingTaskRequestsKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreAllowIncomingTaskRequestsKey, userID)
}

type listStore struct {
	api plugin.API
}

// NewListStore creates a new listStore
func NewListStore(api plugin.API) ListStore {
	return &listStore{
		api: api,
	}
}

func (l *listStore) SaveIssue(issue *Issue) error {
	jsonIssue, jsonErr := json.Marshal(issue)
	if jsonErr != nil {
		return jsonErr
	}

	appErr := l.api.KVSet(issueKey(issue.ID), jsonIssue)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (l *listStore) GetIssue(issueID string) (*Issue, error) {
	originalJSONIssue, appErr := l.api.KVGet(issueKey(issueID))
	if appErr != nil {
		return nil, errors.New(appErr.Error())
	}

	if originalJSONIssue == nil {
		return nil, errors.New("cannot find issue")
	}

	var issue *Issue
	err := json.Unmarshal(originalJSONIssue, &issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (l *listStore) RemoveIssue(issueID string) error {
	appErr := l.api.KVDelete(issueKey(issueID))
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (l *listStore) GetAndRemoveIssue(issueID string) (*Issue, error) {
	issue, err := l.GetIssue(issueID)
	if err != nil {
		return nil, err
	}

	err = l.RemoveIssue(issueID)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (l *listStore) GetIssueReference(userID, issueID, listID string) (*IssueRef, int, error) {
	originalJSONList, err := l.api.KVGet(listKey(userID, listID))
	if err != nil {
		return nil, 0, err
	}

	if originalJSONList == nil {
		return nil, 0, errors.New("cannot load list")
	}

	var list []*IssueRef
	jsonErr := json.Unmarshal(originalJSONList, &list)
	if jsonErr != nil {
		list, _, jsonErr = l.legacyIssueRef(userID, listID)
		if list == nil {
			return nil, 0, jsonErr
		}
	}

	for i, ir := range list {
		if ir.IssueID == issueID {
			return ir, i, nil
		}
	}
	return nil, 0, errors.New("cannot find issue")
}

func (l *listStore) GetIssueListAndReference(userID, issueID string) (string, *IssueRef, int) {
	ir, n, _ := l.GetIssueReference(userID, issueID, MyListKey)
	if ir != nil {
		return MyListKey, ir, n
	}

	ir, n, _ = l.GetIssueReference(userID, issueID, OutListKey)
	if ir != nil {
		return OutListKey, ir, n
	}

	ir, n, _ = l.GetIssueReference(userID, issueID, InListKey)
	if ir != nil {
		return InListKey, ir, n
	}

	return "", nil, 0
}

func (l *listStore) AddReference(userID, issueID, listID, foreignUserID, foreignIssueID string) error {
	for i := 0; i < StoreRetries; i++ {
		list, originalJSONList, err := l.getList(userID, listID)
		if err != nil {
			return err
		}

		for _, ir := range list {
			if ir.IssueID == issueID {
				return errors.New("issue id already exists in list")
			}
		}

		list = append(list, &IssueRef{
			IssueID:        issueID,
			ForeignIssueID: foreignIssueID,
			ForeignUserID:  foreignUserID,
		})

		ok, err := l.saveList(userID, listID, list, originalJSONList)
		if err != nil {
			return err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return nil
		}
	}

	return errors.New("unable to store installation")
}

func (l *listStore) RemoveReference(userID, issueID, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		list, originalJSONList, err := l.getList(userID, listID)
		if err != nil {
			return err
		}

		found := false
		for i, ir := range list {
			if ir.IssueID == issueID {
				list = append(list[:i], list[i+1:]...)
				found = true
			}
		}

		if !found {
			return errors.New("cannot find issue")
		}

		ok, err := l.saveList(userID, listID, list, originalJSONList)
		if err != nil {
			return err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return nil
		}
	}

	return errors.New("unable to store list")
}

func (l *listStore) PopReference(userID, listID string) (*IssueRef, error) {
	for i := 0; i < StoreRetries; i++ {
		list, originalJSONList, err := l.getList(userID, listID)
		if err != nil {
			return nil, err
		}

		if len(list) == 0 {
			return nil, errors.New("cannot find issue")
		}

		ir := list[0]
		list = list[1:]

		ok, err := l.saveList(userID, listID, list, originalJSONList)
		if err != nil {
			return nil, err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return ir, nil
		}
	}

	return nil, errors.New("unable to store list")
}

func (l *listStore) BumpReference(userID, issueID, listID string) error {
	for i := 0; i < StoreRetries; i++ {
		list, originalJSONList, err := l.getList(userID, listID)
		if err != nil {
			return err
		}

		var i int
		var ir *IssueRef

		for i, ir = range list {
			if issueID == ir.IssueID {
				break
			}
		}

		if i == len(list) {
			return errors.New("cannot find issue")
		}

		newList := append([]*IssueRef{ir}, list[:i]...)
		newList = append(newList, list[i+1:]...)

		ok, err := l.saveList(userID, listID, newList, originalJSONList)
		if err != nil {
			return err
		}

		// If err is nil but ok is false, then something else updated the installs between the get and set above
		// so we need to try again, otherwise we can return
		if ok {
			return nil
		}
	}

	return errors.New("unable to store list")
}

func (l *listStore) GetList(userID, listID string) ([]*IssueRef, error) {
	irs, _, err := l.getList(userID, listID)
	return irs, err
}

func (l *listStore) getList(userID, listID string) ([]*IssueRef, []byte, error) {
	originalJSONList, err := l.api.KVGet(listKey(userID, listID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONList == nil {
		return []*IssueRef{}, originalJSONList, nil
	}

	var list []*IssueRef
	jsonErr := json.Unmarshal(originalJSONList, &list)
	if jsonErr != nil {
		return l.legacyIssueRef(userID, listID)
	}

	return list, originalJSONList, nil
}

func (l *listStore) saveList(userID, listID string, list []*IssueRef, originalJSONList []byte) (bool, error) {
	newJSONList, jsonErr := json.Marshal(list)
	if jsonErr != nil {
		return false, jsonErr
	}

	ok, appErr := l.api.KVCompareAndSet(listKey(userID, listID), originalJSONList, newJSONList)
	if appErr != nil {
		return false, errors.New(appErr.Error())
	}

	return ok, nil
}

func (l *listStore) legacyIssueRef(userID, listID string) ([]*IssueRef, []byte, error) {
	originalJSONList, err := l.api.KVGet(listKey(userID, listID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONList == nil {
		return []*IssueRef{}, originalJSONList, nil
	}

	var list []string
	jsonErr := json.Unmarshal(originalJSONList, &list)
	if jsonErr != nil {
		return nil, nil, jsonErr
	}

	newList := []*IssueRef{}
	for _, v := range list {
		newList = append(newList, &IssueRef{IssueID: v})
	}

	return newList, originalJSONList, nil
}

func (p *Plugin) saveLastReminderTimeForUser(userID string) error {
	strTime := strconv.FormatInt(model.GetMillis(), 10)
	appErr := p.API.KVSet(reminderKey(userID), []byte(strTime))
	if appErr != nil {
		return errors.New(appErr.Error())
	}
	return nil
}

func (p *Plugin) getLastReminderTimeForUser(userID string) (int64, error) {
	timeBytes, appErr := p.API.KVGet(reminderKey(userID))
	if appErr != nil {
		return 0, errors.New(appErr.Error())
	}

	if timeBytes == nil {
		return 0, nil
	}

	reminderAt, err := strconv.ParseInt(string(timeBytes), 10, 64)
	if err != nil {
		return 0, err
	}

	return reminderAt, nil
}

func (p *Plugin) saveReminderPreference(userID string, preference bool) error {
	preferenceString := strconv.FormatBool(preference)
	appErr := p.API.KVSet(reminderEnabledKey(userID), []byte(preferenceString))
	if appErr != nil {
		return appErr
	}
	return nil
}

// getReminderPreference - gets user preference on reminder - default value will be true if in case any error
func (p *Plugin) getReminderPreference(userID string) bool {
	preferenceByte, appErr := p.API.KVGet(reminderEnabledKey(userID))
	if appErr != nil {
		p.API.LogError("error getting the reminder preference, err=", appErr.Error())
		return true
	}

	if preferenceByte == nil {
		p.API.LogInfo(`reminder preference is empty. Defaulting to "on"`)
		return true
	}

	preference, err := strconv.ParseBool(string(preferenceByte))
	if err != nil {
		p.API.LogError("unable to parse the reminder preference, err=", err.Error())
		return true
	}

	return preference
}

func (p *Plugin) saveAllowIncomingTaskRequestsPreference(userID string, preference bool) error {
	preferenceString := strconv.FormatBool(preference)
	appErr := p.API.KVSet(allowIncomingTaskRequestsKey(userID), []byte(preferenceString))
	if appErr != nil {
		return appErr
	}
	return nil
}

// getAllowIncomingTaskRequestsPreference - gets user preference on allowing incoming task requests from other users - default value will be true if in case any error
func (p *Plugin) getAllowIncomingTaskRequestsPreference(userID string) (bool, error) {
	preferenceByte, appErr := p.API.KVGet(allowIncomingTaskRequestsKey(userID))
	if appErr != nil {
		err := errors.Wrap(appErr, "error getting the allow incoming task requests preference")
		return true, err
	}

	if preferenceByte == nil {
		p.API.LogDebug(`allow incoming task requests is empty. Defaulting to "on"`)
		return true, nil
	}

	preference, err := strconv.ParseBool(string(preferenceByte))
	if err != nil {
		err := errors.Wrap(appErr, "unable to parse the allow incoming task requests preference")
		return true, err
	}

	return preference, nil
}

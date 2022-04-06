package main

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	fbClient "github.com/mattermost/focalboard/server/client"
	fbModel "github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

const (
	// StoreBoardToUserKey is the key used to map a user ID to a board ID
	StoreUserToBoardKey = "user_to_board"
)

type focalboardListStore struct {
	client *fbClient.Client
	api    plugin.API
}

func userToBoardKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreUserToBoardKey, userID)
}

// NewFocalboardListStore creates a new focalboardListStore
func NewFocalboardListStore(api plugin.API, client *fbClient.Client) FocalboardListStore {
	return &focalboardListStore{
		api:    api,
		client: client,
	}
}

func (l *focalboardListStore) getOrCreateBoardForUser(userID string) (*fbModel.Board, error) {
	rawBoardID, appErr := l.api.KVGet(userToBoardKey(userID))
	if appErr != nil {
		return nil, errors.Wrap(appErr, "unable to get board id from user id")
	}

	if rawBoardID == nil {
		teams, appErr := l.api.GetTeamMembersForUser(userID, 0, 1)
		if appErr != nil {
			return nil, errors.Wrap(appErr, "unable to get team members for user")
		}

		if teams == nil || len(teams) == 0 {
			return nil, errors.New(fmt.Sprintf("user %s not on any teams", userID))
		}

		selfDM, appErr := l.api.GetDirectChannel(userID, userID)
		if appErr != nil {
			return nil, errors.Wrap(appErr, "unable to get self DM for user")
		}

		now := model.GetMillis()

		board := &fbModel.Board{
			ID:         model.NewId(),
			TeamID:     teams[0].TeamId,
			ChannelID:  selfDM.Id,
			Type:       fbModel.BoardTypePrivate,
			Title:      "To Do",
			CreatedBy:  userID,
			Properties: map[string]interface{}{},
			CardProperties: []map[string]interface{}{
				{
					"id":      model.NewId(),
					"name":    "Created By",
					"type":    "createdBy",
					"options": []interface{}{},
				},
				{
					"id":      model.NewId(),
					"name":    "Created At",
					"type":    "createdTime",
					"options": []interface{}{},
				},
				{
					"id":   model.NewId(),
					"name": "Status",
					"type": "select",
					"options": []map[string]interface{}{
						{
							"id":    model.NewId(),
							"value": "Inbox",
							"color": "propColorGray",
						},
						{
							"id":    model.NewId(),
							"value": "To Do",
							"color": "propColorYellow",
						},
						{
							"id":    model.NewId(),
							"value": "Done",
							"color": "propColorGreen",
						},
						{
							"id":    model.NewId(),
							"value": "Won't Do",
							"color": "propColorRed",
						},
					},
				},
			},
			ColumnCalculations: map[string]interface{}{},
			CreateAt:           now,
			UpdateAt:           now,
			DeleteAt:           0,
		}

		block := fbModel.Block{
			ID:       model.NewId(),
			Type:     fbModel.TypeView,
			BoardID:  board.ID,
			ParentID: board.ID,
			Schema:   1,
			Fields: map[string]interface{}{
				"viewType":           fbModel.TypeBoard,
				"sortOptions":        []interface{}{},
				"visiblePropertyIds": []interface{}{},
				"visibleOptionIds":   []interface{}{},
				"hiddenOptionIds":    []interface{}{},
				"collapsedOptionIds": []interface{}{},
				"filter": map[string]interface{}{
					"operation": "and",
					"filters":   []interface{}{},
				},
				"cardOrder":          []interface{}{},
				"columnWidths":       map[string]interface{}{},
				"columnCalculations": map[string]interface{}{},
				"kanbanCalculations": map[string]interface{}{},
				"defaultTemplateId":  "",
			},
			Title:    "All",
			CreateAt: now,
			UpdateAt: now,
			DeleteAt: 0,
		}

		boardsAndBlocks := &fbModel.BoardsAndBlocks{Boards: []*fbModel.Board{board}, Blocks: []fbModel.Block{block}}

		boardsAndBlocks, resp := l.client.CreateBoardsAndBlocks(boardsAndBlocks)
		if resp.Error != nil {
			fmt.Println(resp.StatusCode)
			return nil, errors.Wrap(resp.Error, "unable to create board")
		}
		if len(boardsAndBlocks.Boards) == 0 {
			return nil, errors.New("no board returned")
		}

		board = boardsAndBlocks.Boards[0]

		member := &fbModel.BoardMember{
			BoardID:      board.ID,
			UserID:       userID,
			SchemeEditor: true,
		}

		_, resp = l.client.AddMemberToBoard(member)
		if resp.Error != nil {
			return nil, errors.Wrap(resp.Error, "unable to add user to board")
		}

		appErr = l.api.KVSet(userToBoardKey(userID), []byte(board.ID))
		if appErr != nil {
			return nil, errors.Wrap(appErr, "unable to store board id for user")
		}

		return board, nil
	}

	boardID := string(rawBoardID)
	board, resp := l.client.GetBoard(boardID, "")
	if resp.Error != nil {
		return nil, errors.Wrap(resp.Error, "unable to get board by id")
	}

	return board, nil
}

func getCardPropertyByName(board *fbModel.Board, name string) map[string]interface{} {
	for _, prop := range board.CardProperties {
		if prop["name"] == name {
			return prop
		}
	}

	return nil
}

func getPropertyOptionByValue(property map[string]interface{}, value string) map[string]interface{} {
	optionInterfaces, ok := property["options"].([]interface{})
	if !ok {
		return nil
	}

	for _, optionInterface := range optionInterfaces {
		option, ok := optionInterface.(map[string]interface{})
		if !ok {
			continue
		}

		if option["value"] == value {
			return option
		}
	}

	return nil
}

func (l *focalboardListStore) AddIssue(userID string, issue *Issue) error {
	board, err := l.getOrCreateBoardForUser(userID)
	if err != nil {
		return err
	}

	statusProp := getCardPropertyByName(board, "Status")
	if statusProp == nil {
		return errors.New("status card property not found on board")
	}

	todoOption := getPropertyOptionByValue(statusProp, "To Do")
	if todoOption == nil {
		return errors.New("to do option not found on status card property")
	}

	now := model.GetMillis()

	card := fbModel.Block{
		BoardID:   board.ID,
		Type:      fbModel.TypeCard,
		Title:     issue.Message,
		CreatedBy: userID,
		Fields: map[string]interface{}{
			"icon": "📋",
			"properties": map[string]interface{}{
				statusProp["id"].(string): todoOption["id"],
			},
		},
		CreateAt: now,
		UpdateAt: now,
		DeleteAt: 0,
	}

	_, resp := l.client.InsertBlocks(board.ID, []fbModel.Block{card})
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func (l *focalboardListStore) GetIssue(issueID string) (*Issue, error) {
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

func (l *focalboardListStore) RemoveIssue(issueID string) error {
	appErr := l.api.KVDelete(issueKey(issueID))
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func (l *focalboardListStore) GetAndRemoveIssue(issueID string) (*Issue, error) {
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

func (l *focalboardListStore) GetIssueReference(userID, issueID, listID string) (*IssueRef, int, error) {
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

func (l *focalboardListStore) GetIssueListAndReference(userID, issueID string) (string, *IssueRef, int) {
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

func (l *focalboardListStore) AddReference(userID, issueID, listID, foreignUserID, foreignIssueID string) error {
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

func (l *focalboardListStore) RemoveReference(userID, issueID, listID string) error {
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

func (l *focalboardListStore) PopReference(userID, listID string) (*IssueRef, error) {
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

func (l *focalboardListStore) BumpReference(userID, issueID, listID string) error {
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

func (l *focalboardListStore) GetList(userID, listID string) ([]*IssueRef, error) {
	irs, _, err := l.getList(userID, listID)
	return irs, err
}

func (l *focalboardListStore) getList(userID, listID string) ([]*IssueRef, []byte, error) {
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

func (l *focalboardListStore) saveList(userID, listID string, list []*IssueRef, originalJSONList []byte) (bool, error) {
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

func (l *focalboardListStore) legacyIssueRef(userID, listID string) ([]*IssueRef, []byte, error) {
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

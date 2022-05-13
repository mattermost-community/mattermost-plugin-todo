package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/pkg/errors"

	fbClient "github.com/mattermost/focalboard/server/client"
	fbModel "github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

const (
	// StoreBoardToUserKey is the key used to map a user ID to a board ID
	StoreUserToBoardKey = "user_to_board"

	// StoreReferenceListKey is the key used to store a list of reference issues
	StoreReferenceListKey = "reference_list"

	StatusInbox  = "Inbox"
	StatusToDo   = "To Do"
	StatusDone   = "Done"
	StatusWontDo = "Won't Do"
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

func (l *focalboardListStore) getBoardIDForUser(userID string) (string, error) {
	rawBoardID, appErr := l.api.KVGet(userToBoardKey(userID))
	if appErr != nil {
		return "", errors.Wrap(appErr, "unable to get board id from user id")
	}

	if rawBoardID == nil {
		return "", nil
	}

	return string(rawBoardID), nil
}

func (l *focalboardListStore) getOrCreateBoardForUser(userID string) (*fbModel.Board, error) {
	boardID, err := l.getBoardIDForUser(userID)
	if err != nil {
		return nil, err
	}

	if boardID == "" {
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
					"type":    "person",
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
							"value": StatusInbox,
							"color": "propColorGray",
						},
						{
							"id":    model.NewId(),
							"value": StatusToDo,
							"color": "propColorYellow",
						},
						{
							"id":    model.NewId(),
							"value": StatusDone,
							"color": "propColorGreen",
						},
						{
							"id":    model.NewId(),
							"value": StatusWontDo,
							"color": "propColorRed",
						},
					},
				},
				{
					"id":      model.NewId(),
					"name":    "Post ID",
					"type":    "text",
					"options": []interface{}{},
				},
			},
			CreateAt: now,
			UpdateAt: now,
			DeleteAt: 0,
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
		fmt.Println(resp.StatusCode)
		if boardsAndBlocks == nil {
			return nil, errors.New("no boards or blocks returned")
		}
		if len(boardsAndBlocks.Boards) == 0 {
			return nil, errors.New("no board returned")
		}

		board = boardsAndBlocks.Boards[0]

		member := &fbModel.BoardMember{
			BoardID:     board.ID,
			UserID:      userID,
			SchemeAdmin: true,
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

	board, resp := l.client.GetBoard(boardID, "")
	if resp.Error != nil {
		return nil, errors.Wrap(resp.Error, "unable to get board by id")
	}

	return board, nil
}

func (l *focalboardListStore) getBoardForUser(userID string) (*fbModel.Board, error) {
	boardID, err := l.getBoardIDForUser(userID)
	if err != nil {
		return nil, err
	}

	if boardID == "" {
		return nil, nil
	}

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

func getPropertyValueForCard(block *fbModel.Block, propertyID string) *string {
	if block.Type != fbModel.TypeCard {
		return nil
	}

	properties, ok := block.Fields["properties"].(map[string]interface{})
	if !ok {
		return nil
	}

	value, ok := properties[propertyID].(string)
	if !ok {
		return nil
	}

	return &value
}

func (l *focalboardListStore) AddIssue(userID string, issue *ExtendedIssue) error {
	board, err := l.getOrCreateBoardForUser(userID)
	if err != nil {
		return err
	}

	statusProp := getCardPropertyByName(board, "Status")
	if statusProp == nil {
		return errors.New("status card property not found on board")
	}

	creator := userID
	optionTitle := "To Do"
	if issue.ForeignUser != "" {
		creator = issue.ForeignUser
		optionTitle = StatusInbox
	}
	statusOption := getPropertyOptionByValue(statusProp, optionTitle)
	if statusOption == nil {
		return errors.New("option not found on status card property")
	}

	postIDProp := getCardPropertyByName(board, "Post ID")
	if postIDProp == nil {
		return errors.New("post id card property not found on board")
	}

	createdByProp := getCardPropertyByName(board, "Created By")
	if createdByProp == nil {
		return errors.New("created by card property not found on board")
	}

	now := model.GetMillis()

	card := fbModel.Block{
		BoardID:   board.ID,
		Type:      fbModel.TypeCard,
		Title:     issue.Message,
		CreatedBy: creator,
		Fields: map[string]interface{}{
			"icon": "📋",
			"properties": map[string]interface{}{
				statusProp["id"].(string):    statusOption["id"],
				postIDProp["id"].(string):    issue.PostID,
				createdByProp["id"].(string): creator,
			},
		},
		CreateAt: now,
		UpdateAt: now,
		DeleteAt: 0,
	}

	blocks, resp := l.client.InsertBlocks(board.ID, []fbModel.Block{card})
	if resp.Error != nil {
		return resp.Error
	}

	if len(blocks) != 1 {
		return errors.New("blocks not inserted correctly")
	}

	issue.ID = blocks[0].ID

	return nil
}

func (l *focalboardListStore) GetIssuesByListType(userID string) (map[string][]*ExtendedIssue, error) {
	board, err := l.getBoardForUser(userID)
	if err != nil {
		return nil, err
	}

	lists := map[string][]*ExtendedIssue{
		MyListKey: {},
		InListKey: {},
	}

	if board == nil {
		return lists, nil
	}

	if board == nil {
		return lists, nil
	}

	blocks, resp := l.client.GetAllBlocksForBoard(board.ID)
	if resp.Error != nil {
		return nil, errors.Wrap(resp.Error, "unable to get blocks for board")
	}

	statusProp := getCardPropertyByName(board, "Status")
	if statusProp == nil {
		return nil, errors.New("status card property not found on board")
	}

	todoOption := getPropertyOptionByValue(statusProp, StatusToDo)
	if todoOption == nil {
		return nil, errors.New("to do option not found on status card property")
	}

	inboxOption := getPropertyOptionByValue(statusProp, StatusInbox)
	if inboxOption == nil {
		return nil, errors.New("inbox option not found on status card property")
	}

	var cardOrder []string
	for _, b := range blocks {
		if b.Type == fbModel.TypeView {
			cardOrderInt := b.Fields["cardOrder"].([]interface{})
			cardOrder = make([]string, len(cardOrderInt))
			for index, strInt := range cardOrderInt {
				cardOrder[index] = strInt.(string)
			}
			continue
		}

		status := getPropertyValueForCard(&b, statusProp["id"].(string))
		if status == nil {
			continue
		}

		switch *status {
		case todoOption["id"].(string):
			lists[MyListKey] = append(lists[MyListKey], convertBlockToExtendedIssue(board, &b, userID))
		case inboxOption["id"].(string):
			lists[InListKey] = append(lists[InListKey], convertBlockToExtendedIssue(board, &b, userID))
		}
	}

	fmt.Printf("%v\n", lists[MyListKey])

	if cardOrder != nil {
		sort.Slice(lists[MyListKey], func(i, j int) bool {
			return indexForSorting(cardOrder, lists[MyListKey][i].ID) < indexForSorting(cardOrder, lists[MyListKey][j].ID)
		})
	}

	return lists, nil
}

func indexForSorting(strSlice []string, str string) int {
	for i := range strSlice {
		if strSlice[i] == str {
			return i
		}
	}
	return math.MaxInt
}

func (l *focalboardListStore) GetIssue(userID, issueID string) (*ExtendedIssue, error) {
	board, err := l.getBoardForUser(userID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get board")
	}
	if board == nil {
		return nil, errors.New("unable to find board")
	}

	return l.getIssueFromBoard(board, issueID, userID)
}

func (l *focalboardListStore) getIssueFromBoard(board *fbModel.Board, issueID, userID string) (*ExtendedIssue, error) {
	block, err := l.getBlock(board.ID, issueID)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, nil
	}
	return convertBlockToExtendedIssue(board, block, userID), nil
}

func (l *focalboardListStore) getBlock(boardID, blockID string) (*fbModel.Block, error) {
	blocks, resp := l.client.GetAllBlocksForBoard(boardID)
	if resp.Error != nil {
		return nil, errors.Wrap(resp.Error, "unable to get blocks")
	}

	for _, b := range blocks {
		if b.ID == blockID {
			return &b, nil
		}
	}
	return nil, nil
}

func (l *focalboardListStore) DeleteIssue(userID, issueID string) error {
	board, err := l.getBoardForUser(userID)
	if err != nil {
		return errors.Wrap(err, "unable to get board")
	}
	if board == nil {
		return errors.New("unable to find board")
	}

	issue, err := l.getIssueFromBoard(board, issueID, userID)
	if err != nil {
		return errors.Wrap(err, "unable to get issue")
	}
	if issue == nil {
		return nil
	}

	err = l.RemoveReference(issue.ForeignUser, issue.ID)
	if err != nil {
		l.api.LogDebug("unable to remove reference")
	}

	_, resp := l.client.DeleteBlock(board.ID, issueID)
	if resp.Error != nil {
		return errors.Wrap(resp.Error, "unable to delete block")
	}

	return nil
}

func (l *focalboardListStore) UpdateIssueStatus(userID, issueID, status string) (*ExtendedIssue, string, error) {
	board, err := l.getBoardForUser(userID)
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to get board")
	}
	if board == nil {
		return nil, "", errors.New("unable to find board")
	}

	block, err := l.getBlock(board.ID, issueID)
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to get block to update issue status")
	}
	if block == nil {
		return nil, "", errors.New("unable to find block to update")
	}

	statusProp := getCardPropertyByName(board, "Status")
	if statusProp == nil {
		return nil, "", errors.New("status card property not found on board")
	}
	statusID := statusProp["id"].(string)

	todoOption := getPropertyOptionByValue(statusProp, StatusToDo)
	if todoOption == nil {
		return nil, "", errors.New("to do option not found on status card property")
	}
	todoID := todoOption["id"].(string)

	inboxOption := getPropertyOptionByValue(statusProp, StatusInbox)
	if inboxOption == nil {
		return nil, "", errors.New("inbox option not found on status card property")
	}
	inboxID := inboxOption["id"].(string)

	newOption := getPropertyOptionByValue(statusProp, status)
	if newOption == nil {
		return nil, "", errors.New("new option not found on status card property")
	}
	newID := newOption["id"].(string)

	statusOptionIDPtr := getPropertyValueForCard(block, statusID)
	if statusOptionIDPtr == nil {
		return nil, "", errors.New("unable to find value for status property")
	}

	statusOptionID := *statusOptionIDPtr
	listType := ""

	switch statusOptionID {
	case todoID:
		listType = MyListKey
	case inboxID:
		listType = InListKey
	}

	properties, ok := block.Fields["properties"].(map[string]interface{})
	if !ok {
		return nil, "", errors.New("unable to get block properties")
	}
	properties[statusID] = newID

	patch := &fbModel.BlockPatch{
		UpdatedFields: map[string]interface{}{
			"properties": properties,
		},
	}

	_, resp := l.client.PatchBlock(board.ID, block.ID, patch)
	if resp.Error != nil {
		return nil, "", errors.Wrap(err, "unable to patch block")
	}

	return convertBlockToExtendedIssue(board, block, userID), listType, nil
}

func convertBlockToExtendedIssue(board *fbModel.Board, block *fbModel.Block, userID string) *ExtendedIssue {
	createdByProperty := getCardPropertyByName(board, "Created By")
	if createdByProperty == nil {
		fmt.Println("createdByProperty is nil")
		return nil
	}

	createdByValue := getPropertyValueForCard(block, createdByProperty["id"].(string))
	foreignUserID := ""
	if createdByValue != nil && *createdByValue != userID {
		foreignUserID = *createdByValue
	}

	return &ExtendedIssue{
		Issue: Issue{
			ID:       block.ID,
			Message:  block.Title,
			CreateAt: block.CreateAt,
		},
		ForeignUser: foreignUserID,
		ForeignList: "myfb",
	}
}

func referenceListKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreReferenceListKey, userID)
}

func (l *focalboardListStore) GetReferenceList(userID string) ([]*ExtendedIssue, []byte, error) {
	originalJSONList, err := l.api.KVGet(referenceListKey(userID))
	if err != nil {
		return nil, nil, err
	}

	if originalJSONList == nil {
		return []*ExtendedIssue{}, nil, nil
	}

	var list []*ExtendedIssue
	jsonErr := json.Unmarshal(originalJSONList, &list)
	if jsonErr != nil {
		return nil, nil, errors.Wrap(jsonErr, "unable to decode reference list")
	}

	return list, originalJSONList, nil
}

func (l *focalboardListStore) AddReference(issue *ExtendedIssue) error {
	for i := 0; i < StoreRetries; i++ {
		list, originalJSONList, err := l.GetReferenceList(issue.ForeignUser)
		if err != nil {
			return err
		}

		for _, i := range list {
			if i.ID == issue.ID {
				return errors.New("issue id already exists in list")
			}
		}

		list = append(list, issue)

		ok, err := l.saveReferenceList(issue.ForeignUser, list, originalJSONList)
		if err != nil {
			return err
		}

		if ok {
			return nil
		}
	}

	return errors.New("unable to add reference")
}

func (l *focalboardListStore) saveReferenceList(userID string, list []*ExtendedIssue, originalJSONList []byte) (bool, error) {
	newJSONList, jsonErr := json.Marshal(list)
	if jsonErr != nil {
		return false, errors.Wrap(jsonErr, "unable to encode reference list")
	}

	ok, appErr := l.api.KVCompareAndSet(referenceListKey(userID), originalJSONList, newJSONList)
	if appErr != nil {
		return false, errors.Wrap(appErr, "unable to compare and set referenec list")
	}

	return ok, nil
}

func (l *focalboardListStore) RemoveReference(userID, issueID string) error {
	for i := 0; i < StoreRetries; i++ {
		list, originalJSONList, err := l.GetReferenceList(userID)
		if err != nil {
			return err
		}

		found := false
		for i, ir := range list {
			if ir.ID == issueID {
				list = append(list[:i], list[i+1:]...)
				found = true
			}
		}

		if !found {
			return nil
		}

		ok, err := l.saveReferenceList(userID, list, originalJSONList)
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

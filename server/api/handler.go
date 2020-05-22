package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mattermost/mattermost-plugin-todo/server/command"
	"github.com/mattermost/mattermost-plugin-todo/server/todo"
	"github.com/mattermost/mattermost-server/v5/model"
)

type addAPIRequest struct {
	Message string `json:"message"`
	SendTo  string `json:"send_to"`
	PostID  string `json:"post_id"`
}

func (s *Service) handleAdd(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var addRequest *addAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&addRequest)
	if err != nil {
		s.log.LogError("Unable to decode JSON err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	senderName := s.todo.GetUserName(userID)

	if addRequest.SendTo == "" {
		err = s.todo.AddIssue(userID, addRequest.Message, addRequest.PostID)
		if err != nil {
			s.log.LogError("Unable to add issue err=" + err.Error())
			s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to add issue", err)
			return
		}
		replyMessage := fmt.Sprintf("@%s attached a todo to this thread", senderName)
		s.postReplyIfNeeded(addRequest.PostID, replyMessage, addRequest.Message)
		return
	}

	receiver, appErr := s.plugin.GetUserByUsername(addRequest.SendTo)
	if appErr != nil {
		s.log.LogError("username not valid, err=" + appErr.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to find user", err)
		return
	}

	if receiver.Id == userID {
		err = s.todo.AddIssue(userID, addRequest.Message, addRequest.PostID)
		if err != nil {
			s.log.LogError("Unable to add issue err=" + err.Error())
			s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to add issue", err)
			return
		}
		replyMessage := fmt.Sprintf("@%s attached a todo to this thread", senderName)
		s.postReplyIfNeeded(addRequest.PostID, replyMessage, addRequest.Message)
		return
	}

	issueID, err := s.todo.SendIssue(userID, receiver.Id, addRequest.Message, addRequest.PostID)

	if err != nil {
		s.log.LogError("Unable to send issue err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to send issue", err)
		return
	}

	receiverMessage := fmt.Sprintf("You have received a new Todo from @%s", senderName)
	s.plugin.SendRefreshEvent(receiver.Id)
	s.bot.PostBotCustomDM(receiver.Id, receiverMessage, addRequest.Message, issueID)

	replyMessage := fmt.Sprintf("@%s sent @%s a todo attached to this thread", senderName, addRequest.SendTo)
	s.postReplyIfNeeded(addRequest.PostID, replyMessage, addRequest.Message)
}

func (s *Service) handleList(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	listInput := r.URL.Query().Get("list")
	listID := todo.MyListKey
	switch listInput {
	case "out":
		listID = todo.OutListKey
	case "in":
		listID = todo.InListKey
	}

	issues, err := s.todo.GetIssueList(userID, listID)
	if err != nil {
		s.log.LogError("Unable to get issues for user err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to get issues for user", err)
		return
	}

	if len(issues.Issues) > 0 && r.URL.Query().Get("reminder") == "true" {
		var lastReminderAt int64
		lastReminderAt, err = s.todo.GetReminder(userID)
		if err != nil {
			s.log.LogError("Unable to send reminder err=" + err.Error())
			s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to send reminder", err)
			return
		}

		var timezone *time.Location
		offset, _ := strconv.Atoi(r.Header.Get("X-Timezone-Offset"))
		timezone = time.FixedZone("local", -60*offset)

		// Post reminder message if it's the next day and been more than an hour since the last post
		now := model.GetMillis()
		nt := time.Unix(now/1000, 0).In(timezone)
		lt := time.Unix(lastReminderAt/1000, 0).In(timezone)
		if nt.Sub(lt).Hours() >= 1 && (nt.Day() != lt.Day() || nt.Month() != lt.Month() || nt.Year() != lt.Year()) {
			s.bot.PostBotDM(userID, "Daily Reminder:\n\n"+issues.ToString())
			s.todo.AddReminder(userID)
		}
	}

	issuesJSON, err := json.Marshal(issues.Issues)
	if err != nil {
		s.log.LogError("Unable marhsal issues list to json err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable marhsal issues list to json", err)
		return
	}

	w.Write(issuesJSON)
}

type acceptAPIRequest struct {
	ID string `json:"id"`
}

func (s *Service) handleAccept(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var acceptRequest *acceptAPIRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&acceptRequest); err != nil {
		s.log.LogError("Unable to decode JSON err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	todoMessage, sender, err := s.todo.AcceptIssue(userID, acceptRequest.ID)

	if err != nil {
		s.log.LogError("Unable to accept issue err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to accept issue", err)
		return
	}

	userName := s.todo.GetUserName(userID)

	message := fmt.Sprintf("@%s accepted a Todo you sent: %s", userName, todoMessage)
	s.plugin.SendRefreshEvent(sender)
	s.bot.PostBotDM(sender, message)
}

type completeAPIRequest struct {
	ID string `json:"id"`
}

func (s *Service) handleComplete(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var completeRequest *completeAPIRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&completeRequest); err != nil {
		s.log.LogError("Unable to decode JSON err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	issue, foreignID, err := s.todo.CompleteIssue(userID, completeRequest.ID)
	if err != nil {
		s.log.LogError("Unable to complete issue err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to complete issue", err)
		return
	}

	userName := s.todo.GetUserName(userID)
	replyMessage := fmt.Sprintf("@%s completed a todo attached to this thread", userName)
	s.postReplyIfNeeded(issue.PostID, replyMessage, issue.Message)

	if foreignID == "" {
		return
	}

	message := fmt.Sprintf("@%s completed a Todo you sent: %s", userName, issue.Message)
	s.plugin.SendRefreshEvent(foreignID)
	s.bot.PostBotDM(foreignID, message)
}

type removeAPIRequest struct {
	ID string `json:"id"`
}

func (s *Service) handleRemove(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var removeRequest *removeAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&removeRequest)
	if err != nil {
		s.log.LogError("Unable to decode JSON err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	issue, foreignID, isSender, err := s.todo.RemoveIssue(userID, removeRequest.ID)
	if err != nil {
		s.log.LogError("Unable to remove issue, err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to remove issue", err)
		return
	}

	userName := s.todo.GetUserName(userID)
	replyMessage := fmt.Sprintf("@%s removed a todo attached to this thread", userName)
	s.postReplyIfNeeded(issue.PostID, replyMessage, issue.Message)

	if foreignID == "" {
		return
	}

	message := fmt.Sprintf("@%s removed a Todo you received: %s", userName, issue.Message)
	if isSender {
		message = fmt.Sprintf("@%s declined a Todo you sent: %s", userName, issue.Message)
	}

	s.plugin.SendRefreshEvent(foreignID)
	s.bot.PostBotDM(foreignID, message)
}

type bumpAPIRequest struct {
	ID string `json:"id"`
}

func (s *Service) handleBump(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var bumpRequest *bumpAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bumpRequest)
	if err != nil {
		s.log.LogError("Unable to decode JSON err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	todoMessage, foreignUser, foreignIssueID, err := s.todo.BumpIssue(userID, bumpRequest.ID)
	if err != nil {
		s.log.LogError("Unable to bump issue, err=" + err.Error())
		s.handleErrorWithCode(w, http.StatusInternalServerError, "Unable to bump issue", err)
		return
	}

	if foreignUser == "" {
		return
	}

	userName := s.todo.GetUserName(userID)

	message := fmt.Sprintf("@%s bumped a Todo you received.", userName)

	s.plugin.SendRefreshEvent(foreignUser)
	s.bot.PostBotCustomDM(foreignUser, message, todoMessage, foreignIssueID)
}

func (s *Service) handleErrorWithCode(w http.ResponseWriter, code int, errTitle string, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(struct {
		Error   string `json:"error"`
		Details string `json:"details"`
	}{
		Error:   errTitle,
		Details: err.Error(),
	})
	_, _ = w.Write(b)
}

func (s *Service) postReplyIfNeeded(postID, message, todo string) {
	if postID != "" {
		err := s.bot.ReplyPostBot(postID, message, todo)
		if err != nil {
			s.log.LogError(err.Error())
		}
	}
}

func (s *Service) executeCommand(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	args := model.CommandArgsFromJson(r.Body)

	com := command.NewCommand(args, s.log, s.plugin, s.todo)
	responses := com.Handle()

	for _, response := range responses {
		if response.Type != command.ResponseTypePayload {
			continue
		}

		w.Write(response.ToJSON())
		return
	}
}

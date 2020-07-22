package main

const (
	sourceCommand = "command"
	sourceWebapp  = "webapp"
)

func (p *Plugin) trackCommand(userID, command string) {
	p.tracker.TrackUserEvent("command", userID, map[string]interface{}{
		"command": command,
	})
}

func (p *Plugin) trackAddIssue(userID, source string, attached bool) {
	p.tracker.TrackUserEvent("add_issue", userID, map[string]interface{}{
		"source":   source,
		"attached": attached,
	})
}

func (p *Plugin) trackSendIssue(userID, source string, attached bool) {
	p.tracker.TrackUserEvent("send_issue", userID, map[string]interface{}{
		"source":   source,
		"attached": attached,
	})
}

func (p *Plugin) trackCompleteIssue(userID string) {
	p.tracker.TrackUserEvent("complete_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackRemoveIssue(userID string) {
	p.tracker.TrackUserEvent("remove_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackAcceptIssue(userID string) {
	p.tracker.TrackUserEvent("accept_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackBumpIssue(userID string) {
	p.tracker.TrackUserEvent("bump_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackFrontend(userID, event string, properties map[string]interface{}) {
	if properties == nil {
		properties = map[string]interface{}{}
	}
	p.tracker.TrackUserEvent("frontend_"+event, userID, properties)
}

func (p *Plugin) trackDailySummary(userID string) {
	p.tracker.TrackUserEvent("daily_summary_sent", userID, map[string]interface{}{})
}

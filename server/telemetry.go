package main

type telemetrySource string

const (
	sourceCommand telemetrySource = "command"
	sourceWebapp  telemetrySource = "webapp"
)

func (p *Plugin) trackCommand(userID, command string) {
	_ = p.tracker.TrackUserEvent("command", userID, map[string]interface{}{
		"command": command,
	})
}

func (p *Plugin) trackAddIssue(userID string, source telemetrySource, attached bool) {
	_ = p.tracker.TrackUserEvent("add_issue", userID, map[string]interface{}{
		"source":   source,
		"attached": attached,
	})
}

func (p *Plugin) trackSendIssue(userID string, source telemetrySource, attached bool) {
	_ = p.tracker.TrackUserEvent("send_issue", userID, map[string]interface{}{
		"source":   source,
		"attached": attached,
	})
}

func (p *Plugin) trackCompleteIssue(userID string) {
	_ = p.tracker.TrackUserEvent("complete_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackRemoveIssue(userID string) {
	_ = p.tracker.TrackUserEvent("remove_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackAcceptIssue(userID string) {
	_ = p.tracker.TrackUserEvent("accept_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackEditIssue(userID string) {
	_ = p.tracker.TrackUserEvent("edit_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackChangeAssignment(userID string) {
	_ = p.tracker.TrackUserEvent("change_issue_assignment", userID, map[string]interface{}{})
}

func (p *Plugin) trackBumpIssue(userID string) {
	_ = p.tracker.TrackUserEvent("bump_issue", userID, map[string]interface{}{})
}

func (p *Plugin) trackFrontend(userID, event string, properties map[string]interface{}) {
	if properties == nil {
		properties = map[string]interface{}{}
	}
	_ = p.tracker.TrackUserEvent("frontend_"+event, userID, properties)
}

func (p *Plugin) trackDailySummary(userID string) {
	_ = p.tracker.TrackUserEvent("daily_summary_sent", userID, map[string]interface{}{})
}

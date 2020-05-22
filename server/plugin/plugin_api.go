package plugin

import "github.com/mattermost/mattermost-server/v5/model"

const (
	// WSEventRefresh is the WebSocket event for refreshing the Todo list
	WSEventRefresh = "refresh"
)

func (p *Plugin) SendRefreshEvent(userID string) {
	p.API.PublishWebSocketEvent(
		WSEventRefresh,
		nil,
		&model.WebsocketBroadcast{UserId: userID},
	)
}

func (p *Plugin) GetUserByUsername(name string) (*model.User, *model.AppError) {
	return p.API.GetUserByUsername(name)
}

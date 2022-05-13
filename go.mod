module github.com/mattermost/mattermost-plugin-todo

go 1.16

replace github.com/mattermost/focalboard/server => /usr/local/joram/go/src/github.com/mattermost/focalboard/server

require (
	github.com/mattermost/focalboard/server v0.0.0-20220512192951-9dbb0c88e162
	github.com/mattermost/mattermost-plugin-api v0.0.27
	github.com/mattermost/mattermost-server/v6 v6.5.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.1
)

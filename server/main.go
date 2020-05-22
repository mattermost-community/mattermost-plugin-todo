package main

import (
	todoplugin "github.com/mattermost/mattermost-plugin-todo/server/plugin"
	mattermost "github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	mattermost.ClientMain(&todoplugin.Plugin{})
}

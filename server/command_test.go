package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/mock"
)

func TestSetttingsCommand(t *testing.T) {
	api := &plugintest.API{}
	api.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	api.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)

	apiKVSetFailed := &plugintest.API{}
	apiKVSetFailed.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	apiKVSetFailed.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(model.NewAppError("failed", "", nil, "", 400))

	tests := []struct {
		name    string
		api     *plugintest.API
		args    []string
		wantErr bool
		want    bool
	}{
		{
			name:    "Setting summary successful",
			api:     api,
			args:    []string{"summary", "on"},
			wantErr: false,
			want:    false,
		},
		{
			name:    "Setting summary failed due to KVSet failed",
			api:     apiKVSetFailed,
			args:    []string{"summary", "on"},
			wantErr: true,
			want:    false,
		},
		{
			name:    "Setting summary failed due to invalid number of arguments",
			api:     api,
			args:    []string{"summary", "on", "extraValue"},
			wantErr: true,
			want:    true,
		},
		{
			name:    "Setting summary failed due to no arguments",
			api:     api,
			args:    []string{},
			wantErr: true,
			want:    true,
		},
		{
			name:    "Setting summary failed due to invalid argument",
			api:     api,
			args:    []string{"summary", "test"},
			wantErr: true,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := Plugin{}
			plugin.SetAPI(tt.api)

			resp, err := plugin.runSettingsCommand(tt.args, &model.CommandArgs{})
			if tt.wantErr != (err != nil) {
				t.Errorf("runSettingsCommand wantErr= %v got err= %v", tt.wantErr, err)
				return
			}

			if tt.want != resp {
				t.Errorf("runSettingsCommand got= %v want= %v", resp, tt.want)
			}
		})
	}
}

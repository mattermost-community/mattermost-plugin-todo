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
	api.On("KVGet", mock.AnythingOfType("string"), mock.Anything).Return([]byte("true"), nil)

	apiKVSetFailed := &plugintest.API{}
	apiKVSetFailed.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	apiKVSetFailed.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(model.NewAppError("failed", "", nil, "", 400))
	apiKVSetFailed.On("LogDebug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))

	tests := []struct {
		name    string
		api     *plugintest.API
		args    []string
		wantErr bool
		want    bool
	}{
		{
			name:    "Setting successful without any arguments",
			api:     api,
			args:    []string{},
			wantErr: false,
			want:    false,
		},
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
			name:    "Setting summary successful due to no arguments",
			api:     api,
			args:    []string{"summary"},
			wantErr: false,
			want:    false,
		},
		{
			name:    "Setting summary failed due to invalid argument",
			api:     api,
			args:    []string{"summary", "test"},
			wantErr: true,
			want:    true,
		},
		{
			name:    "Setting allow_incoming_task_requests successful",
			api:     api,
			args:    []string{"allow_incoming_task_requests", "on"},
			wantErr: false,
			want:    false,
		},
		{
			name:    "Setting allow_incoming_task_requests failed due to KVSet failed",
			api:     apiKVSetFailed,
			args:    []string{"allow_incoming_task_requests", "on"},
			wantErr: true,
			want:    false,
		},
		{
			name:    "Setting allow_incoming_task_requests failed due to invalid number of arguments",
			api:     api,
			args:    []string{"allow_incoming_task_requests", "on", "extraValue"},
			wantErr: true,
			want:    true,
		},
		{
			name:    "Setting allow_incoming_task_requests successful due to no arguments",
			api:     api,
			args:    []string{"allow_incoming_task_requests"},
			wantErr: false,
			want:    false,
		},
		{
			name:    "Setting allow_incoming_task_requests failed due to invalid argument",
			api:     api,
			args:    []string{"allow_incoming_task_requests", "test"},
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

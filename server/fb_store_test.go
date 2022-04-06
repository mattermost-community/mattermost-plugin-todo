package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPropertyOptionByValue(t *testing.T) {
	property := map[string]interface{}{
		"id":   "1",
		"name": "Status",
		"type": "select",
		"options": []map[string]interface{}{
			{
				"id":    "2",
				"value": "Inbox",
				"color": "propColorGray",
			},
			{
				"id":    "3",
				"value": "To Do",
				"color": "propColorYellow",
			},
			{
				"id":    "4",
				"value": "Done",
				"color": "propColorGreen",
			},
			{
				"id":    "4",
				"value": "Won't Do",
				"color": "propColorRed",
			},
		},
	}

	option := getPropertyOptionByValue(property, "To Do")
	require.NotNil(t, option)
	assert.Equal(t, "To Do", option["value"])

	option = getPropertyOptionByValue(property, "Junk")
	assert.Nil(t, option)
}

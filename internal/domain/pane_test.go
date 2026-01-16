package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestPaneJSONMarshaling(t *testing.T) {
	pane := Pane{
		Path:    "/home/user",
		Command: "vim",
		Split:   "vertical",
		Size:    50,
	}

	data, err := json.Marshal(pane)
	assert.NoError(t, err)

	var unmarshaled Pane
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, pane, unmarshaled)
}

func TestPaneYAMLMarshaling(t *testing.T) {
	pane := Pane{
		Path:    "/home/user",
		Command: "vim",
		Split:   "horizontal",
		Size:    30,
	}

	data, err := yaml.Marshal(pane)
	assert.NoError(t, err)

	var unmarshaled Pane
	err = yaml.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, pane, unmarshaled)
}

func TestWindowWithPanes(t *testing.T) {
	window := Window{
		Name: "editor",
		Path: "/home/user/project",
		Panes: []Pane{
			{Command: "vim", Split: "vertical", Size: 50},
			{Command: "npm run dev", Split: "horizontal", Size: 30},
		},
	}

	data, err := json.Marshal(window)
	assert.NoError(t, err)

	var unmarshaled Window
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, window.Name, unmarshaled.Name)
	assert.Len(t, unmarshaled.Panes, 2)
	assert.Equal(t, "vim", unmarshaled.Panes[0].Command)
}

func TestWindowWithoutPanes(t *testing.T) {
	window := Window{
		Name:    "editor",
		Path:    "/home/user/project",
		Command: "vim",
	}

	data, err := json.Marshal(window)
	assert.NoError(t, err)

	var unmarshaled Window
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, window.Name, unmarshaled.Name)
	assert.Empty(t, unmarshaled.Panes)
}

package domain

type Pane struct {
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
	Command string `json:"command,omitempty" yaml:"command,omitempty"`
	Split   string `json:"split,omitempty" yaml:"split,omitempty"` // "horizontal" or "vertical"
	Size    int    `json:"size,omitempty" yaml:"size,omitempty"`   // percentage
}

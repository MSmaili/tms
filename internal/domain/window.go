package domain

type Window struct {
	Name    string `json:"name" yaml:"name"`
	Path    string `json:"path" yaml:"path"`
	Index   *int   `json:"index,omitempty" yaml:"index,omitempty"`
	Layout  string `json:"layout,omitempty" yaml:"layout,omitempty"`
	Command string `json:"command,omitempty" yaml:"command,omitempty"`
	Panes   []Pane `json:"panes,omitempty" yaml:"panes,omitempty"`
}

package tmux

type Window struct {
	Name   string
	Path   string
	Layout string
	Panes  []Pane
}

type Pane struct {
	Path    string
	Command string
}

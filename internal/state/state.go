package state

// State represents tmux state for diffing purposes
type State struct {
	Sessions map[string]*Session
}

type Session struct {
	Name    string
	Windows []*Window
}

type Window struct {
	Name  string
	Path  string
	Panes []*Pane
}

type Pane struct {
	Index   int
	Path    string
	Command string
}

func NewState() *State {
	return &State{Sessions: make(map[string]*Session)}
}

func (s *State) AddSession(name string) *Session {
	session := &Session{Name: name}
	s.Sessions[name] = session
	return session
}

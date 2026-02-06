package tmux

import (
	"fmt"
	"strings"

	"github.com/MSmaili/muxie/internal/backend"
	"github.com/MSmaili/muxie/internal/plan"
)

type TmuxBackend struct {
	client        Client
	paneBaseIndex int
}

func init() {
	backend.Register("tmux", func() (backend.Backend, error) {
		return NewBackend()
	})
}

func NewBackend() (*TmuxBackend, error) {
	c, err := New()
	if err != nil {
		return nil, err
	}
	return &TmuxBackend{client: c}, nil
}

func (b *TmuxBackend) Name() string {
	return "tmux"
}

func (b *TmuxBackend) QueryState() (backend.StateResult, error) {
	result, err := RunQuery(b.client, LoadStateQuery{})
	if err != nil {
		return backend.StateResult{}, err
	}

	b.paneBaseIndex = result.PaneBaseIndex

	sessions := make([]backend.Session, len(result.Sessions))
	for i, s := range result.Sessions {
		windows := make([]backend.Window, len(s.Windows))
		for j, w := range s.Windows {
			panes := make([]backend.Pane, len(w.Panes))
			for k, p := range w.Panes {
				panes[k] = backend.Pane{
					Index:   k,
					Path:    p.Path,
					Command: p.Command,
				}
			}
			windows[j] = backend.Window{
				Name:   w.Name,
				Path:   w.Path,
				Layout: w.Layout,
				Panes:  panes,
			}
		}
		sessions[i] = backend.Session{
			Name:    s.Name,
			Windows: windows,
		}
	}

	return backend.StateResult{
		Sessions: sessions,
		Active: backend.ActiveContext{
			Session: result.Active.Session,
			Window:  result.Active.Window,
			Pane:    result.Active.Pane,
			Path:    result.Active.Path,
		},
	}, nil
}

func (b *TmuxBackend) Apply(actions []backend.Action) error {
	tmuxActions := b.mapActions(actions)
	return b.client.ExecuteBatch(tmuxActions)
}

func (b *TmuxBackend) DryRun(actions []backend.Action) []string {
	tmuxActions := b.mapActions(actions)
	lines := make([]string, len(tmuxActions))
	for i, a := range tmuxActions {
		lines[i] = "tmux " + strings.Join(a.Args(), " ")
	}
	return lines
}

func (b *TmuxBackend) Attach(session string) error {
	return b.client.Attach(session)
}

func (b *TmuxBackend) mapActions(actions []backend.Action) []Action {
	result := make([]Action, 0, len(actions))
	for _, a := range actions {
		if ta := b.mapAction(a); ta != nil {
			result = append(result, ta)
		}
	}
	return result
}

func (b *TmuxBackend) mapAction(a backend.Action) Action {
	switch action := a.(type) {
	case plan.CreateSessionAction:
		return CreateSession{Name: action.Name, WindowName: action.WindowName, Path: action.Path}
	case plan.CreateWindowAction:
		return CreateWindow{Session: action.Session, Name: action.Name, Path: action.Path}
	case plan.SplitPaneAction:
		return SplitPane{Target: fmt.Sprintf("%s:%s", action.Session, action.Window), Path: action.Path}
	case plan.SendKeysAction:
		return SendKeys{Target: fmt.Sprintf("%s:%s.%d", action.Session, action.Window, action.Pane+b.paneBaseIndex), Keys: action.Command}
	case plan.SelectLayoutAction:
		return SelectLayout{Target: fmt.Sprintf("%s:%s", action.Session, action.Window), Layout: action.Layout}
	case plan.ZoomPaneAction:
		return ZoomPane{Target: fmt.Sprintf("%s:%s.%d", action.Session, action.Window, action.Pane+b.paneBaseIndex)}
	case plan.KillSessionAction:
		return KillSession{Name: action.Name}
	case plan.KillWindowAction:
		return KillWindow{Target: fmt.Sprintf("%s:%s", action.Session, action.Window)}
	default:
		return nil
	}
}

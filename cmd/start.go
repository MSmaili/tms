package cmd

import (
	"fmt"

	"github.com/MSmaili/tms/internal/manifest"
	"github.com/MSmaili/tms/internal/plan"
	"github.com/MSmaili/tms/internal/state"
	"github.com/MSmaili/tms/internal/tmux"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [workspace-name-or-path]",
	Short: "Start a tmux workspace",
	RunE:  runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	var nameOrPath string
	if len(args) > 0 {
		nameOrPath = args[0]
	}

	resolver := manifest.NewResolver()
	workspacePath, err := resolver.Resolve(nameOrPath)
	if err != nil {
		return err
	}

	loader := manifest.NewFileLoader(workspacePath)
	workspace, err := loader.Load()
	if err != nil {
		return fmt.Errorf("loading workspace: %w", err)
	}

	client, err := tmux.New()
	if err != nil {
		return fmt.Errorf("initializing tmux client: %w", err)
	}

	desired := manifestToState(workspace)

	actual, err := queryTmuxState(client)
	if err != nil {
		return fmt.Errorf("querying tmux state: %w", err)
	}

	stateDiff := state.Compare(desired, actual)

	planDiff := stateDiffToPlanDiff(stateDiff, desired)
	strategy := &plan.MergeStrategy{}
	p := strategy.Plan(planDiff)

	if p.IsEmpty() {
		fmt.Println("Workspace already up to date")
	} else {
		actions := planActionsToTmuxActions(p.Actions)
		for _, action := range actions {
			if err := client.Execute(action); err != nil {
				return fmt.Errorf("executing action: %w", err)
			}
		}
	}

	for sessionName := range workspace.Sessions {
		return client.Attach(sessionName)
	}

	return nil
}

func manifestToState(ws *manifest.Workspace) *state.State {
	s := state.NewState()
	for sessionName, windows := range ws.Sessions {
		session := s.AddSession(sessionName)
		for _, w := range windows {
			window := &state.Window{Name: w.Name, Path: w.Path}
			for _, p := range w.Panes {
				window.Panes = append(window.Panes, &state.Pane{Path: p.Path, Command: p.Command})
			}
			session.Windows = append(session.Windows, window)
		}
	}
	return s
}

func queryTmuxState(client tmux.Client) (*state.State, error) {
	s := state.NewState()

	sessions, err := tmux.RunQuery(client, tmux.ListSessionsQuery{})
	if err != nil {
		return s, nil
	}

	for _, sess := range sessions {
		session := s.AddSession(sess.Name)

		windows, err := tmux.RunQuery(client, tmux.ListWindowsQuery{Session: sess.Name})
		if err != nil {
			continue
		}

		for _, w := range windows {
			window := &state.Window{Name: w.Name, Path: w.Path}

			panes, err := tmux.RunQuery(client, tmux.ListPanesQuery{Target: fmt.Sprintf("%s:%s", sess.Name, w.Name)})
			if err == nil {
				for _, p := range panes {
					window.Panes = append(window.Panes, &state.Pane{Path: p.Path, Command: p.Command})
				}
			}

			session.Windows = append(session.Windows, window)
		}
	}

	return s, nil
}

func stateDiffToPlanDiff(sd state.Diff, desired *state.State) plan.Diff {
	pd := plan.Diff{
		Windows: make(map[string]plan.ItemDiff[plan.Window]),
	}

	for _, sessionName := range sd.Sessions.Missing {
		session := desired.Sessions[sessionName]
		ps := plan.Session{Name: sessionName}
		for _, w := range session.Windows {
			pw := plan.Window{Name: w.Name, Path: w.Path}
			for _, p := range w.Panes {
				pw.Panes = append(pw.Panes, plan.Pane{Path: p.Path, Command: p.Command})
			}
			ps.Windows = append(ps.Windows, pw)
		}
		pd.Sessions.Missing = append(pd.Sessions.Missing, ps)
	}

	for _, sessionName := range sd.Sessions.Extra {
		pd.Sessions.Extra = append(pd.Sessions.Extra, plan.Session{Name: sessionName})
	}

	for sessionName, wd := range sd.Windows {
		pwd := plan.ItemDiff[plan.Window]{}

		for _, w := range wd.Missing {
			pw := plan.Window{Name: w.Name, Path: w.Path}
			for _, p := range w.Panes {
				pw.Panes = append(pw.Panes, plan.Pane{Path: p.Path, Command: p.Command})
			}
			pwd.Missing = append(pwd.Missing, pw)
		}

		for _, w := range wd.Extra {
			pwd.Extra = append(pwd.Extra, plan.Window{Name: w.Name, Path: w.Path})
		}

		for _, m := range wd.Mismatched {
			pwd.Mismatched = append(pwd.Mismatched, plan.Mismatch[plan.Window]{
				Desired: stateWindowToPlanWindow(m.Desired),
				Actual:  stateWindowToPlanWindow(m.Actual),
			})
		}

		pd.Windows[sessionName] = pwd
	}

	return pd
}

func stateWindowToPlanWindow(w state.Window) plan.Window {
	pw := plan.Window{Name: w.Name, Path: w.Path}
	for _, p := range w.Panes {
		pw.Panes = append(pw.Panes, plan.Pane{Path: p.Path, Command: p.Command})
	}
	return pw
}

func planActionsToTmuxActions(actions []plan.Action) []tmux.Action {
	var result []tmux.Action
	for _, a := range actions {
		result = append(result, planActionToTmuxAction(a))
	}
	return result
}

func planActionToTmuxAction(a plan.Action) tmux.Action {
	switch action := a.(type) {
	case plan.CreateSessionAction:
		return tmux.CreateSession{Name: action.Name, Path: action.Path}
	case plan.CreateWindowAction:
		return tmux.CreateWindow{Session: action.Session, Name: action.Name, Path: action.Path}
	case plan.SplitPaneAction:
		return tmux.SplitPane{Target: action.Target, Path: action.Path}
	case plan.SendKeysAction:
		return tmux.SendKeys{Target: action.Target, Keys: action.Command}
	case plan.KillSessionAction:
		return tmux.KillSession{Name: action.Name}
	case plan.KillWindowAction:
		return tmux.KillWindow{Target: action.Target}
	default:
		return nil
	}
}

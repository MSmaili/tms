package plan

import "fmt"

type Strategy interface {
	Plan(diff Diff) *Plan
}

type MergeStrategy struct{}

func (s *MergeStrategy) Plan(diff Diff) *Plan {
	plan := &Plan{Actions: []Action{}}
	createMissing(plan, diff)
	return plan
}

type ForceStrategy struct{}

func (s *ForceStrategy) Plan(diff Diff) *Plan {
	plan := &Plan{Actions: []Action{}}
	killExtra(plan, diff)
	recreateMismatched(plan, diff)
	createMissing(plan, diff)
	return plan
}

func killExtra(plan *Plan, diff Diff) {
	for _, session := range diff.Sessions.Extra {
		plan.Actions = append(plan.Actions, KillSessionAction{Name: session.Name})
	}

	for sessionName, windowDiff := range diff.Windows {
		for _, window := range windowDiff.Extra {
			plan.Actions = append(plan.Actions, KillWindowAction{
				Target: fmt.Sprintf("%s:%s", sessionName, window.Name),
			})
		}
	}
}

func recreateMismatched(plan *Plan, diff Diff) {
	for sessionName, windowDiff := range diff.Windows {
		for _, mismatch := range windowDiff.Mismatched {
			plan.Actions = append(plan.Actions, KillWindowAction{
				Target: fmt.Sprintf("%s:%s", sessionName, mismatch.Actual.Name),
			})
			createWindow(plan, sessionName, mismatch.Desired)
		}
	}
}

func createMissing(plan *Plan, diff Diff) {
	for _, session := range diff.Sessions.Missing {
		createSession(plan, session)
	}

	for sessionName, windowDiff := range diff.Windows {
		for _, window := range windowDiff.Missing {
			createWindow(plan, sessionName, window)
		}
	}
}

func createSession(plan *Plan, session Session) {
	if len(session.Windows) == 0 {
		return
	}

	firstWindow := session.Windows[0]
	plan.Actions = append(plan.Actions, CreateSessionAction{
		Name: session.Name,
		Path: firstWindow.Path,
	})
	addPanesAndCommands(plan, session.Name, firstWindow)

	for _, window := range session.Windows[1:] {
		createWindow(plan, session.Name, window)
	}
}

func createWindow(plan *Plan, sessionName string, window Window) {
	plan.Actions = append(plan.Actions, CreateWindowAction{
		Session: sessionName,
		Name:    window.Name,
		Path:    window.Path,
	})
	addPanesAndCommands(plan, sessionName, window)
}

func addPanesAndCommands(plan *Plan, sessionName string, window Window) {
	if len(window.Panes) == 0 {
		return
	}

	target := fmt.Sprintf("%s:%s", sessionName, window.Name)

	for _, pane := range window.Panes[1:] {
		plan.Actions = append(plan.Actions, SplitPaneAction{
			Target: target,
			Path:   pane.Path,
		})
	}

	for i, pane := range window.Panes {
		if pane.Command != "" {
			plan.Actions = append(plan.Actions, SendKeysAction{
				Target:  fmt.Sprintf("%s.%d", target, i),
				Command: pane.Command,
			})
		}
	}
}

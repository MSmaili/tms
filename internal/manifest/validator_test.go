package manifest

import (
	"strings"
	"testing"
)

func TestValidate_ValidWorkspace(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{
				Name: "dev",
				Windows: []Window{
					{Name: "editor", Path: "/home/user"},
					{Name: "terminal", Path: "/home/user"},
				},
			},
		},
	}

	errs := Validate(ws)
	if len(errs) > 0 {
		t.Errorf("expected valid workspace to pass validation, got: %v", errs)
	}
}

func TestValidate_NilWorkspace(t *testing.T) {
	errs := Validate(nil)
	if len(errs) == 0 {
		t.Error("expected error for nil workspace")
	}
	if !strings.Contains(errs[0].Message, "nil") {
		t.Errorf("expected 'nil' in error, got: %v", errs)
	}
}

func TestValidate_EmptySessions(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{},
	}

	errs := Validate(ws)
	if len(errs) == 0 {
		t.Error("expected error for empty sessions")
	}
	if !strings.Contains(errs[0].Message, "no sessions") {
		t.Errorf("expected 'no sessions' in error, got: %v", errs)
	}
}

func TestValidate_EmptySessionName(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{
				Name:    "",
				Windows: []Window{{Name: "editor", Path: "/home"}},
			},
		},
	}

	errs := Validate(ws)
	if len(errs) == 0 {
		t.Error("expected error for empty session name")
	}
	if !strings.Contains(errs[0].Message, "session name cannot be empty") {
		t.Errorf("expected 'session name cannot be empty' in error, got: %v", errs)
	}
}

func TestValidate_EmptyWindowList(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{Name: "dev", Windows: []Window{}},
		},
	}

	errs := Validate(ws)
	if len(errs) == 0 {
		t.Error("expected error for empty window list")
	}
	if !strings.Contains(errs[0].Message, "no windows") {
		t.Errorf("expected 'no windows' in error, got: %v", errs)
	}
}

func TestValidate_DuplicateWindowNames(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{
				Name: "dev",
				Windows: []Window{
					{Name: "editor", Path: "/home"},
					{Name: "editor", Path: "/home"},
				},
			},
		},
	}

	errs := Validate(ws)
	if len(errs) == 0 {
		t.Error("expected error for duplicate window names")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "duplicate window name") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'duplicate window name' in errors, got: %v", errs)
	}
}

func TestValidate_DuplicateSessionNames(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{Name: "dev", Windows: []Window{{Name: "a", Path: "/home"}}},
			{Name: "dev", Windows: []Window{{Name: "b", Path: "/home"}}},
		},
	}

	errs := Validate(ws)
	if len(errs) == 0 {
		t.Error("expected error for duplicate session names")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "duplicate session name") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'duplicate session name' in errors, got: %v", errs)
	}
}

func TestValidate_WindowsWithoutNames(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{
				Name: "dev",
				Windows: []Window{
					{Path: "/home/user"},
					{Path: "/home/other"},
				},
			},
		},
	}

	errs := Validate(ws)
	if len(errs) > 0 {
		t.Errorf("expected windows without names to be valid, got: %v", errs)
	}
}

func TestValidate_MultipleZoomedPanes(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{
				Name: "dev",
				Windows: []Window{
					{
						Name: "editor",
						Path: "/home",
						Panes: []Pane{
							{Path: "/home/user", Zoom: true},
							{Path: "/home/user", Zoom: true},
						},
					},
				},
			},
		},
	}

	errs := Validate(ws)
	if len(errs) == 0 {
		t.Error("expected error for multiple zoomed panes in one window")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "zoom=true") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'zoom=true' in errors, got: %v", errs)
	}
}

func TestValidate_SingleZoomedPane(t *testing.T) {
	ws := &Workspace{
		Sessions: []Session{
			{
				Name: "dev",
				Windows: []Window{
					{
						Name: "editor",
						Path: "/home",
						Panes: []Pane{
							{Path: "/home/user", Zoom: true},
							{Path: "/home/user"},
						},
					},
				},
			},
		},
	}

	errs := Validate(ws)
	if len(errs) > 0 {
		t.Errorf("expected single zoomed pane to be valid, got: %v", errs)
	}
}

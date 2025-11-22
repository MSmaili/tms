package domain

type Diff struct {
	MissingSessions []string
	ExtraSessions   []string
	CommonSessions  []string

	MissingWindows map[string][]Window
	ExtraWindows   map[string][]Window
	Mismatched     map[string][]WindowMismatch
}

type WindowMismatch struct {
	Desired Window
	Actual  Window
}

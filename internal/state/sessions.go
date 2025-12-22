package state

func compareSessions(diff *Diff, desired, actual *State) {
	for name := range desired.Sessions {
		if _, ok := actual.Sessions[name]; !ok {
			diff.Sessions.Missing = append(diff.Sessions.Missing, name)
		}
	}

	for name := range actual.Sessions {
		if _, ok := desired.Sessions[name]; !ok {
			diff.Sessions.Extra = append(diff.Sessions.Extra, name)
		}
	}
}

func CommonSessions(desired, actual *State) []string {
	common := make([]string, 0, len(desired.Sessions))
	for name := range desired.Sessions {
		if _, exists := actual.Sessions[name]; exists {
			common = append(common, name)
		}
	}
	return common
}

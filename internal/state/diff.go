package state

type Diff struct {
	Sessions ItemDiff[string]
	Windows  map[string]ItemDiff[Window] // key: session name
}

type ItemDiff[T any] struct {
	Missing    []T
	Extra      []T
	Mismatched []Mismatch[T]
}

type Mismatch[T any] struct {
	Desired T
	Actual  T
}

func (d ItemDiff[T]) IsEmpty() bool {
	return len(d.Missing) == 0 && len(d.Extra) == 0 && len(d.Mismatched) == 0
}

func Compare(desired, actual *State) Diff {
	diff := Diff{
		Windows: make(map[string]ItemDiff[Window]),
	}

	compareSessions(&diff, desired, actual)
	compareWindows(&diff, desired, actual)

	return diff
}

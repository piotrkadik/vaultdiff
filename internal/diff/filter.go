package diff

// FilterOptions controls which changes are included in a filtered result set.
type FilterOptions struct {
	// Types is a list of change types to include. If empty, all types are included.
	Types []ChangeType
	// PathPrefix limits results to keys that begin with the given prefix.
	PathPrefix string
}

// Filter returns a new slice of Change values that match the given FilterOptions.
// If opts.Types is empty and opts.PathPrefix is empty, the original slice is
// returned unchanged.
func Filter(changes []Change, opts FilterOptions) []Change {
	if len(opts.Types) == 0 && opts.PathPrefix == "" {
		return changes
	}

	typeSet := make(map[ChangeType]struct{}, len(opts.Types))
	for _, t := range opts.Types {
		typeSet[t] = struct{}{}
	}

	result := make([]Change, 0, len(changes))
	for _, c := range changes {
		if len(typeSet) > 0 {
			if _, ok := typeSet[c.Type]; !ok {
				continue
			}
		}
		if opts.PathPrefix != "" && !hasPrefix(c.Key, opts.PathPrefix) {
			continue
		}
		result = append(result, c)
	}
	return result
}

// hasPrefix reports whether s begins with prefix.
func hasPrefix(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return s[:len(prefix)] == prefix
}

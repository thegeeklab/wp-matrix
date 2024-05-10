package plugin

import "strings"

// EnsurePrefix ensures that the given input string starts with the provided prefix.
func EnsurePrefix(prefix, input string) string {
	if strings.TrimSpace(input) == "" {
		return input
	}

	if strings.HasPrefix(input, prefix) {
		return input
	}

	return prefix + input
}

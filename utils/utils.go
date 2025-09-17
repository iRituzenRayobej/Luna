package utils

import (
	"path/filepath"
)

func FileMatchesPattern(filename, pattern string) bool {
	matched, _ := filepath.Match(pattern, filepath.Base(filename))
	if matched {
		return true
	}
	fullPattern := filepath.FromSlash(pattern)
	matched, _ = filepath.Match(fullPattern, filename)
	return matched
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

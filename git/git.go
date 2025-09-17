package git

import (
	"fmt"
	"luna/config"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetStagedFiles() ([]string, error) {
	output, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return nil, fmt.Errorf("error running git diff --cached: %v", err)
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

func GetFileDiff(filename string) (string, error) {
	diff, err := exec.Command("git", "diff", "--cached", "--", filename).Output()
	if err != nil {
		return "", fmt.Errorf("error getting diff for %s: %v", filename, err)
	}
	return string(diff), nil
}

func ShouldIgnoreFile(filename string, cfg config.Config) bool {
	for _, ignoredFile := range cfg.IgnoredFiles {
		if filename == ignoredFile {
			return true
		}
		if strings.HasSuffix(ignoredFile, "/") && strings.HasPrefix(filename, ignoredFile) {
			return true
		}
	}

	for _, pattern := range cfg.IgnoredPatterns {
		matched, _ := filepath.Match(pattern, filepath.Base(filename))
		if matched {
			return true
		}
		fullPattern := filepath.FromSlash(pattern)
		matched, _ = filepath.Match(fullPattern, filename)
		if matched {
			return true
		}
	}

	return false
}

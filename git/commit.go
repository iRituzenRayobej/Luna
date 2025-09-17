package git

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"luna/ai"
	"luna/config"
)

var emojis = []string{"✨", "🛠️", "🐛", "🔥", "📝", "🚀", "🔧", "🎨", "🔒", "💄"}

func GenerateCommitMessage(apiKey, diff, filename string, cfg config.Config, includeEmoji bool) string {
	commitMsg := ai.CallGemini(apiKey, diff)
	if commitMsg == "" {
		commitMsg = "update " + filename
	}

	hasPrefix := false
	for _, p := range cfg.CommitPrefixes {
		if strings.HasPrefix(strings.ToLower(commitMsg), strings.ToLower(p)) {
			hasPrefix = true
			break
		}
	}

	if !hasPrefix && len(cfg.CommitPrefixes) > 0 {
		rand.Seed(time.Now().UnixNano())
		prefix := cfg.CommitPrefixes[rand.Intn(len(cfg.CommitPrefixes))]
		commitMsg = prefix + " " + commitMsg
	}

	if includeEmoji {
		rand.Seed(time.Now().UnixNano())
		emoji := emojis[rand.Intn(len(emojis))]
		commitMsg = fmt.Sprintf("%s %s", emoji, commitMsg)
	}

	if len(commitMsg) > cfg.MaxCommitLength {
		commitMsg = commitMsg[:cfg.MaxCommitLength-3] + "..."
	}

	return commitMsg
}

func CommitFile(filename, commitMsg string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", commitMsg, "--", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("error committing %s: %v", filename, err)
	}
	return string(out), nil
}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"


type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type RequestBody struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: LunaHelp to see the commands")
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "lunahelp":
		showHelp()
	case "lunacommit":
		runCommitGenerator()
	case "lunaapikey":
		setApiKey()
	default:
		fmt.Println("Unknown command. Use: LunaHelp")
	}
}

func showHelp() {
	fmt.Println("\nLuna - AI Commit Generator")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  LunaHelp       -> Show this help screen")
	fmt.Println("  LunaCommit     -> Generate commit messages for each staged file and commit automatically")
	fmt.Println("  LunaApikey     -> Set your Gemini API key")
}

func runCommitGenerator() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: Set GEMINI_API_KEY using LunaApikey first")
		return
	}

	status, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		fmt.Println("Error running git status:", err)
		return
	}

	lines := strings.Split(string(status), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		file := parts[1]

		// Ignorar arquivos binários
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".exe" || ext == ".dll" || ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			continue
		}

		diff, err := exec.Command("git", "diff", "--cached", "--", file).Output()
		if err != nil || len(diff) == 0 {
			continue
		}

		commitMsg := callGemini(apiKey, string(diff))
		if commitMsg == "" {
			commitMsg = "chore: update " + file
		} else {
			// Garante prefixo tipo "chore:" ou "refactor:" se não tiver
			prefixes := []string{"chore:", "refactor:", "feat:", "fix:", "docs:", "test:"}
			hasPrefix := false
			for _, p := range prefixes {
				if strings.HasPrefix(strings.ToLower(commitMsg), p) {
					hasPrefix = true
					break
				}
			}
			if !hasPrefix {
				commitMsg = "chore: " + commitMsg
			}
		}

		// Faz o commit automático
		cmd := exec.Command("git", "commit", "-m", commitMsg, "--", file)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error committing %s: %s\n", file, string(out))
		} else {
			fmt.Printf("Committed %s with message: %s\n", file, commitMsg)
		}
	}
}

func callGemini(apiKey, diff string) string {
	body := RequestBody{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: fmt.Sprintf("Generate a short commit message for the following diff (include type like chore:, refactor:, feat:, fix:, docs:, test:):\n%s", diff)},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", API_URL+"?key="+apiKey, bytes.NewReader(jsonData))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return ""
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(response.Candidates[0].Content.Parts[0].Text)
	}
	return ""
}

func setApiKey() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: LunaApikey YOUR_API_KEY")
		return
	}

	apiKey := os.Args[2]

	cmd := exec.Command("setx", "GEMINI_API_KEY", apiKey)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error setting API key:", err)
		return
	}

	fmt.Println("GEMINI_API_KEY set successfully! Close and reopen the terminal.")
}

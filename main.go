package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"


const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

var emojis = []string{"✨", "🛠️", "🐛", "🔥", "📝", "🚀", "🔧", "🎨", "🔒", "💄"}

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
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println(Red + "Use: LunaHelp to see the commands" + Reset)
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
		fmt.Println(Red + "Unknown command. Use: LunaHelp" + Reset)
	}
}

func showHelp() {
	fmt.Println(Cyan + `
 __                                     
/  |                                    
$$ |       __    __  _______    ______  
$$ |      /  |  /  |/       \  /      \ 
$$ |      $$ |  $$ |$$$$$$$  | $$$$$$  |
$$ |      $$ |  $$ |$$ |  $$ | /    $$ |
$$ |_____ $$ \__$$ |$$ |  $$ |/$$$$$$$ |
$$       |$$    $$/ $$ |  $$ |$$    $$ |
$$$$$$$$/  $$$$$$/  $$/   $$/  $$$$$$$/ 
                                        
                                        
                                        
made by hax & dan
version: 1.1 (Beta)
` + Reset)

	fmt.Println(Yellow + "Available commands:" + Reset)
	fmt.Println(Green + "  LunaHelp       -> Show this help screen with ASCII art")
	fmt.Println("  LunaCommit     -> Generate commit messages for each staged file with emojis")
	fmt.Println("  LunaApikey     -> Set your Gemini API key" + Reset)
}

func runCommitGenerator() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println(Red + "Error: Set GEMINI_API_KEY using LunaApikey first" + Reset)
		return
	}

	status, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		fmt.Println(Red+"Error running git status:", err, Reset)
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

		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".exe" || ext == ".dll" || ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			fmt.Println(Yellow + "Skipping binary file: " + file + Reset)
			continue
		}

		fmt.Println(Cyan + "Generating commit for file: " + file + Reset)

		diff, err := exec.Command("git", "diff", "--cached", "--", file).Output()
		if err != nil || len(diff) == 0 {
			fmt.Println(Yellow + "No staged changes to commit for file: " + file + Reset)
			continue
		}

		commitMsg := callGemini(apiKey, string(diff))
		if commitMsg == "" {
			commitMsg = "update " + file
		}

		prefixes := []string{"chore:", "refactor:", "feat:", "fix:", "docs:", "test:"}
		hasPrefix := false
		for _, p := range prefixes {
			if strings.HasPrefix(strings.ToLower(commitMsg), p) {
				hasPrefix = true
				break
			}
		}

		if !hasPrefix {
			prefix := prefixes[rand.Intn(len(prefixes))]
			commitMsg = prefix + " " + commitMsg
		}

		emoji := emojis[rand.Intn(len(emojis))]

		fullMsg := fmt.Sprintf("%s %s", emoji, commitMsg)
		if len(fullMsg) > 100 {
			fullMsg = fullMsg[:97] + "..."
		}

		cmd := exec.Command("git", "commit", "-m", fullMsg, "--", file)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf(Red+"Error committing %s: %s\n"+Reset, file, string(out))
		} else {
			fmt.Printf(Green+"Committed %s with message:\n%s\n\n"+Reset, file, fullMsg)
		}
	}
}

func callGemini(apiKey, diff string) string {
	body := RequestBody{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: fmt.Sprintf(
						"Generate a short, concise, one-line commit message for the following diff. "+
							"Keep it under 60 characters, include optional emojis and type like chore:, refactor:, feat:, fix:, docs:, test:\n%s", diff)},
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
		fmt.Println(Red + "Usage: LunaApikey YOUR_API_KEY" + Reset)
		return
	}

	apiKey := os.Args[2]

	cmd := exec.Command("setx", "GEMINI_API_KEY", apiKey)
	err := cmd.Run()
	if err != nil {
		fmt.Println(Red+"Error setting API key:", err, Reset)
		return
	}

	fmt.Println(Green + "GEMINI_API_KEY set successfully! Close and reopen the terminal." + Reset)
}

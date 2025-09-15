package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"

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

	switch os.Args[1] {
	case "LunaHelp":
		showHelp()
	case "LunaCommit":
		runCommitGenerator()
	case "LunaApikey":
		setApiKey()
	default:
		fmt.Println("Unknown command. Use: LunaHelp")
	}
}

func showHelp() {
	fmt.Println("\nLuna - AI Commit Generator")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  LunaHelp       -> Show this help screen")
	fmt.Println("  LunaCommit     -> Generate commit messages for each staged file")
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

		diff, _ := exec.Command("git", "diff", "--cached", "--", file).Output()
		if len(diff) == 0 {
			continue
		}

		commitMsg := callGemini(apiKey, string(diff))
		fmt.Printf("File: %s\nSuggested commit: %s\n\n", file, commitMsg)
	}
}

func callGemini(apiKey, diff string) string {
	body := RequestBody{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: fmt.Sprintf("Generate a short, clear, and descriptive commit message for the following diff:\n%s", diff)},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", API_URL+"?key="+apiKey, bytes.NewReader(jsonData))
	if err != nil {
		return "Error creating request"
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Request error"
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return "Error parsing JSON"
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return response.Candidates[0].Content.Parts[0].Text
	}
	return "Could not generate commit"
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

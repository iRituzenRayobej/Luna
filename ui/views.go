package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"luna/git"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
)

func ShowHelp() {
	fmt.Println(titleStyle.Render("Luna - AI Git Assistant"))
	fmt.Println(`
Available commands:

  LunaHelp (lh, -lh)
  -> Shows this help screen

  LunaCommit (lc, -lc)
  -> Generates commit messages using Gemini AI
  -> Use -e flag for emojis

  LunaApikey (lkey, -lkey)
  -> Sets your Gemini API key

  LunaConfig (config, -config)
  -> Manage project configuration
`)
}

func ShowConfigHelp() {
	fmt.Println(titleStyle.Render("Config Management"))
	fmt.Println(`
  luna config init    - Create project config
  luna config show    - Show current config
  luna config edit    - Edit project config

Config Priority:
  • API Key: Global > Project > Default
  • Other settings: Project > Default
`)
}

func loadFiles() tea.Msg {
	files, err := git.GetStagedFiles()
	return loadedFilesMsg{files: files, err: err}
}

func handleReviewInput(m CommitUI, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		if m.commitIndex >= len(m.files) {
			m.state = stateComplete
		
			return m, func() tea.Msg {
				time.Sleep(2 * time.Second)
				return tea.QuitMsg{}
			}
		}

		filename := m.files[m.commitIndex]
		commitMsg := m.commitMsgs[filename]

		out, err := git.CommitFile(filename, commitMsg)
		if err != nil {
			m.state = stateError
			m.err = err
			return m, nil
		}

		m.commitResults[filename] = out
		m.commitIndex++

		if m.commitIndex >= len(m.files) {
			m.state = stateComplete
			return m, func() tea.Msg {
				time.Sleep(2 * time.Second)
				return tea.QuitMsg{}
			}
		}

		return m, nil

	case "r":
		if m.currentFile >= len(m.files) {
			return m, nil
		}

		filename := m.files[m.currentFile]
		diff, err := git.GetFileDiff(filename)
		if err != nil {
			m.state = stateError
			m.err = err
			return m, nil
		}

		newMsg := git.GenerateCommitMessage(m.cfg.ApiKey, diff, filename, m.cfg, m.includeEmoji)
		m.commitMsgs[filename] = newMsg

		return m, nil

	case "q":
		return m, tea.Quit
	}

	return m, nil
}

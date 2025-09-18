package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"luna/config"
	"luna/git"
)

type CommitUI struct {
	cfg           config.Config
	includeEmoji  bool
	files         []string
	currentFile   int
	commitIndex   int
	commitMsgs    map[string]string
	commitResults map[string]string
	spinner       spinner.Model
	progress      progress.Model
	viewport      viewport.Model
	state         uiState
	err           error
}

type uiState int

const (
	stateLoading uiState = iota
	stateProcessing
	stateReview
	stateComplete
	stateError
)

type fileProcessedMsg struct {
	filename string
	result   string
	err      error
}

type loadedFilesMsg struct {
	files []string
	err   error
}

func InitializeCommitUI(cfg config.Config, includeEmoji bool) CommitUI {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	prog := progress.New(progress.WithDefaultGradient())
	vp := viewport.New(80, 20)

	return CommitUI{
		cfg:           cfg,
		includeEmoji:  includeEmoji,
		commitMsgs:    make(map[string]string),
		commitResults: make(map[string]string),
		spinner:       sp,
		progress:      prog,
		viewport:      vp,
		state:         stateLoading,
	}
}

func (m CommitUI) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadFiles,
	)
}

func (m CommitUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if m.state == stateReview || m.state == stateComplete {
			return handleReviewInput(m, msg)
		}

	case loadedFilesMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}
		if len(msg.files) == 0 {
			m.state = stateComplete
			return m, tea.Quit
		}
		m.files = msg.files
		m.state = stateProcessing
		return m, processNextFile(m)

	case fileProcessedMsg:
		if msg.err != nil {
			m.commitResults[msg.filename] = "Error: " + msg.err.Error()
		} else {
			m.commitResults[msg.filename] = msg.result
		}
		m.currentFile++

		if m.currentFile >= len(m.files) {
			m.state = stateReview
			return m, nil
		}
		return m, processNextFile(m)
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func processNextFile(m CommitUI) tea.Cmd {
	if m.currentFile >= len(m.files) {
		return nil
	}

	file := m.files[m.currentFile]

	return func() tea.Msg {
		if git.ShouldIgnoreFile(file, m.cfg) {
			return fileProcessedMsg{filename: file, result: "Ignored", err: nil}
		}

		diff, err := git.GetFileDiff(file)
		if err != nil {
			return fileProcessedMsg{filename: file, result: "", err: err}
		}

		commitMsg := git.GenerateCommitMessage(m.cfg.ApiKey, diff, file, m.cfg, m.includeEmoji)
		m.commitMsgs[file] = commitMsg

		return fileProcessedMsg{
			filename: file,
			result:   "Commit message generated",
			err:      nil,
		}
	}
}

func (m CommitUI) View() string {
	var b strings.Builder

	switch m.state {
	case stateLoading:
		b.WriteString(fmt.Sprintf("%s Loading files...\n", m.spinner.View()))

	case stateProcessing:
		progressVal := float64(m.currentFile) / float64(len(m.files))
		b.WriteString(fmt.Sprintf("%s Processing files...\n", m.spinner.View()))
		b.WriteString(fmt.Sprintf("Progress: %s\n", m.progress.ViewAs(progressVal)))

		if m.currentFile < len(m.files) {
			b.WriteString(fmt.Sprintf("Current: %s\n", m.files[m.currentFile]))
		}

	case stateReview:
		b.WriteString("📝 Review Commit Messages:\n\n")
		for _, file := range m.files {
			msg := m.commitMsgs[file]
			result := m.commitResults[file]

			status := "✅"
			if strings.Contains(result, "Error") {
				status = "❌"
			}

			b.WriteString(fmt.Sprintf("%s %s: %s\n", status, file, msg))
			b.WriteString(fmt.Sprintf("   Result: %s\n\n", result))
		}
		b.WriteString("Press 'c' to confirm, 'r' to retry, 'q' to quit\n")

	case stateComplete:
		b.WriteString("🎉 All files committed successfully!\n")
		b.WriteString("Exiting in 2 seconds...\n")

	case stateError:
		b.WriteString(fmt.Sprintf("❌ Error: %v\n", m.err))
	}

	return b.String()
}

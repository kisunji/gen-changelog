package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type scene int

const (
	scenePrSelect = iota
	sceneTypeSelect
	sceneBodyInput
)

type model struct {
	scene      scene
	list       list.Model
	body       textinput.Model
	targetPR   string
	targetType string
	fileName   string
	written    bool
	errMsg     string
}

func newModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = ""

	return model{
		list: list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		body: ti,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		// early exit if we're not at top-level directory
		if _, err := os.ReadDir(".changelog"); err != nil {
			return errors.New(".changelog directory not found. Are you in the project root?")
		}

		// -X=GET required since passing -F defaults the request to POST
		cmd := exec.Command("gh", "api", "-X=GET", "search/issues",
			"-F", "q=repo:{owner}/{repo} is:open type:pr head:{branch}",
			"-q", ".items.[].number", // jq filtering
		)

		raw, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("gh error: %s: %v", bytes.TrimSpace(raw), err.Error())
		}

		if len(raw) == 0 {
			return errors.New("no open pull requests found for current branch")
		}

		prs := strings.Split(strings.TrimSpace(string(raw)), "\n")
		// Fast path since most branches will have at most one open PR.
		if len(prs) == 1 {
			return prs[0]
		}

		prListItems := make([]list.Item, len(prs))
		for i, pr := range prs {
			prListItems[i] = item(pr)
		}

		l := list.New(prListItems, itemDelegate{}, 30, len(prListItems)+5)
		l.Title = "Select PR number:"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.SetShowHelp(false)
		l.Styles.Title = lipgloss.NewStyle()
		return l
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case error:
		m.errMsg = msg.Error()
		return m, tea.Quit

	case string: // single PR found
		m.targetPR = msg

	case list.Model: // multiple PRs found
		m.list = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			m.errMsg = "terminated"
			return m, tea.Quit

		case "enter":
			switch m.scene {
			case scenePrSelect:
				if i, ok := m.list.SelectedItem().(item); ok {
					m.targetPR = string(i)
				}
			case sceneTypeSelect:
				if i, ok := m.list.SelectedItem().(item); ok {
					m.targetType = string(i)
				}
			case sceneBodyInput:
				var bb bytes.Buffer
				bb.WriteString(header)
				bb.WriteString(m.targetType)
				bb.WriteString("\n")
				bb.WriteString(m.body.Value())
				bb.WriteString("\n")
				bb.WriteString(footer)

				if err := os.WriteFile(m.fileName, bb.Bytes(), 0644); err != nil {
					m.errMsg = err.Error()
					return m, tea.Quit
				}

				m.written = true
				return m, tea.Quit
			}
		}
	}

	// If targetPR is set, prompt for changelog type
	if m.scene < sceneTypeSelect && m.targetPR != "" {
		m.fileName = fmt.Sprintf(".changelog/%s.txt", m.targetPR)

		if _, err := os.ReadFile(m.fileName); err == nil {
			m.errMsg = m.fileName + " already exists"
			return m, tea.Quit
		}

		m.scene = sceneTypeSelect
		m.list = pickChangelog()
		return m, nil
	}

	// if targetType is set, prompt for body
	if m.scene < sceneBodyInput && m.targetType != "" {
		m.scene = sceneBodyInput
		m.list.SetItems(nil)
		return m, nil
	}

	// handle key presses for subcomponents
	var cmd tea.Cmd
	if m.scene == sceneTypeSelect || m.scene == scenePrSelect {
		m.list, cmd = m.list.Update(msg)
	}
	if m.scene == sceneBodyInput {
		m.body, cmd = m.body.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	var sb strings.Builder

	if m.targetPR != "" {
		sb.WriteString(lipgloss.NewStyle().Bold(true).
			Render(".changelog/"+m.targetPR+".txt") + "\n")
	}

	if m.targetType != "" {
		sb.WriteString(lipgloss.NewStyle().Italic(true).
			Render(header+m.targetType) + "\n")
	}

	if len(m.list.Items()) > 0 {
		sb.WriteString(m.list.View() + "\n")
	}

	if m.scene == sceneBodyInput {
		if m.written {
			sb.WriteString(m.body.Value() + "\n")
			sb.WriteString(footer + "\n")
		} else {
			sb.WriteString(m.body.View() + "\n")
			sb.WriteString(footer + "\n")
			sb.WriteString("\n")
			sb.WriteString(lipgloss.NewStyle().Render("(enter to save or ctrl+c to exit)"))
		}

	}

	if m.written {
		sb.WriteString("File written. \n")
	}

	if m.errMsg != "" {
		sb.WriteString(m.errMsg + "\n")
	}

	if sb.Len() > 0 {
		return sb.String()
	}

	return "loading..."
}

func pickChangelog() list.Model {
	l := list.New(releaseNoteTypes, itemDelegate{}, 30, len(releaseNoteTypes)+5)
	l.Title = "Select a changelog type:"
	l.Styles.Title = lipgloss.NewStyle()
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	return l
}

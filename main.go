package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	header = "```release-note:"
	footer = "```"
)

var releaseNoteTypes = []list.Item{
	item("bug"),
	item("improvement"),
	item("feature"),
	item("security"),
	item("breaking-change"),
	item("deprecation"),
	item("note"),
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

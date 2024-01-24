package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/techierishi/pal/logr"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	paddingLeft         = noStyle.Copy().PaddingLeft(1)
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	inputStyle = paddingLeft.Copy().Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#EEEEEE"})

	focusedButton = focusedStyle.Copy().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Render(" Show ")
	blurredButton = blurredStyle.Copy().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"}).
			Render(" Show ")

	logger = logr.GetLogInstance()
)

type inputModel struct {
	credString string
	hiddenText string
	input      textinput.Model
}

func initialModel(hiddenText string) inputModel {
	m := inputModel{}

	mask := "**********"
	t := textinput.New()
	t.Cursor.Style = focusedStyle
	t.Cursor.TextStyle = focusedStyle
	t.Prompt = ""
	t.PromptStyle = inputStyle
	t.CharLimit = 500
	t.Placeholder = "Password..."
	t.SetValue(mask)
	t.Focus()
	m.input = t
	m.hiddenText = hiddenText

	return m
}

func (m inputModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			s := msg.String()
			if s == "enter" {
				m.input.SetValue(m.hiddenText)
				return m, nil
			}
		}
	}
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *inputModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.credString = m.input.Value()
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m inputModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(m.input.View())
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(helpStyle.Render("[ This system does not support clipboard API]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ press `shift` to show value and then copy manually]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ `esc` to exit ]"))
	b.WriteRune('\n')

	return b.String()
}

func PasswordModal(hiddenText string) *string {
	m := initialModel(hiddenText)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start PasswordModal: %s\n", err)
		os.Exit(1)
	}

	return &m.hiddenText
}

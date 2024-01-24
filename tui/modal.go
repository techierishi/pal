package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	extCode = -1
)

type inputModel struct {
	credString string
	hiddenText string
	focusIndex int
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
			extCode = 0
			return m, tea.Quit
		case "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == 1 {
				m.input.SetValue(m.hiddenText)
				return m, nil
			}

		case "tab", "shift+tab":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 1
			}

			cmds := make([]tea.Cmd, 1)

			if 0 == m.focusIndex {
				cmds[0] = m.input.Focus()
				m.input.TextStyle = focusedStyle
				m.input.Cursor.Style = focusedStyle

			} else if 1 == m.focusIndex {
				m.blur([]string{"input"})
			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *inputModel) blur(inputKeys []string) {
	for _, inputKey := range inputKeys {
		if inputKey == "input" {
			m.input.Blur()
			m.input.Cursor.Style = noStyle
			m.input.Cursor.TextStyle = noStyle
		}
	}
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

	button := &blurredButton
	if m.focusIndex == 1 {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("[ `tab` to switch, `ctrl+c` to cancel ]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ `enter` on button to show value and then copy manually]"))
	b.WriteRune('\n')

	return b.String()
}

func Modal(hiddenText string) *string {

	m := initialModel(hiddenText)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	if extCode == 0 {
		os.Exit(extCode)
	}

	return &m.hiddenText
}

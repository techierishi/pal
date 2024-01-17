package aliasm

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
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
			Render(" Submit ")
	blurredButton = blurredStyle.Copy().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"}).
			Render(" Submit ")

	aliasInfo = AliasInfo{
		Alias:   "",
		Command: "",
	}

	exitCode = -1
)

type aliasModel struct {
	focusIndex int
	alias      textinput.Model
	command    textarea.Model
}

func initialModel() aliasModel {
	m := aliasModel{}

	t := textinput.New()
	t.Cursor.Style = focusedStyle
	t.Cursor.TextStyle = focusedStyle
	t.Prompt = ""
	t.PromptStyle = inputStyle
	t.CharLimit = 500
	t.Placeholder = "Alias Name..."
	t.Focus()
	m.alias = t

	ta := textarea.New()
	ta.Placeholder = "Alias Value..."
	ta.Cursor.Style = blurredStyle
	ta.Cursor.TextStyle = blurredStyle
	ta.ShowLineNumbers = false
	m.command = ta

	return m
}

func (m aliasModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m aliasModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			exitCode = 0
			return m, tea.Quit
		case "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == 2 {
				return m, tea.Quit
			}

		case "tab", "shift+tab":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}

			cmds := make([]tea.Cmd, 2)

			if 0 == m.focusIndex {
				cmds[0] = m.alias.Focus()
				m.alias.TextStyle = focusedStyle
				m.alias.Cursor.Style = focusedStyle

				m.blur([]string{"command"})

			} else if 1 == m.focusIndex {
				cmds[1] = m.command.Focus()
				m.command.Cursor.Style = focusedStyle
				m.command.Cursor.TextStyle = focusedStyle
				m.command.FocusedStyle.Text = focusedStyle

				m.blur([]string{"alias"})

			} else if 2 == m.focusIndex {
				m.blur([]string{"alias", "command"})

			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *aliasModel) blur(inputKeys []string) {

	for _, inputKey := range inputKeys {
		if inputKey == "command" {
			m.command.Blur()
			m.command.Cursor.Style = noStyle
			m.command.Cursor.TextStyle = noStyle
		}
		if inputKey == "alias" {
			m.alias.Blur()
			m.alias.TextStyle = noStyle
			m.alias.Cursor.Style = noStyle
		}
	}

}

func (m *aliasModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.command, cmd = m.command.Update(msg)
	aliasInfo.Command = m.command.Value()
	cmds = append(cmds, cmd)
	m.alias, cmd = m.alias.Update(msg)
	aliasInfo.Alias = m.alias.Value()
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m aliasModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(m.alias.View())
	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString(m.command.View())
	b.WriteRune('\n')

	button := &blurredButton
	if m.focusIndex == 2 {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("[ `tab` to switch, `ctrl+c` to cancel ]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ `enter` on button to save ]"))
	b.WriteRune('\n')

	return b.String()
}

func NewAlias() *AliasInfo {

	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	if exitCode == 0 {
		os.Exit(exitCode)
	}

	return &aliasInfo
}

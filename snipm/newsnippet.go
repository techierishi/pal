package snipm

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

	snippetInfo = SnippetInfo{
		Description: "",
		Command:     "",
		PromptPass:  "",
	}

	exitCode = -1
)

type snippetModel struct {
	focusIndex  int
	snippet     textarea.Model
	description textinput.Model
	promptPass  textinput.Model
}

func initialModel() snippetModel {
	m := snippetModel{}
	ta := textarea.New()
	ta.Placeholder = "Snippet..."
	ta.Cursor.Style = focusedStyle
	ta.Cursor.TextStyle = focusedStyle
	ta.ShowLineNumbers = false
	ta.Focus()
	m.snippet = ta

	var t textinput.Model
	t = textinput.New()
	t.Cursor.Style = blurredStyle
	t.Cursor.TextStyle = blurredStyle
	t.Prompt = ""
	t.PromptStyle = inputStyle
	t.CharLimit = 500
	t.Placeholder = "Description..."
	m.description = t

	t = textinput.New()
	t.Cursor.Style = blurredStyle
	t.Cursor.TextStyle = blurredStyle
	t.Prompt = ""
	t.PromptStyle = inputStyle
	t.CharLimit = 4
	t.Placeholder = "Prompt for password [Y/n], Default = N"
	m.promptPass = t

	return m
}

func (m snippetModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m snippetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			exitCode = 0
			return m, tea.Quit
		case "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == 3 {
				return m, tea.Quit
			}

		case "tab", "shift+tab":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 3 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 3
			}

			cmds := make([]tea.Cmd, 3)

			if 0 == m.focusIndex {
				cmds[0] = m.snippet.Focus()
				m.snippet.Cursor.Style = focusedStyle
				m.snippet.Cursor.TextStyle = focusedStyle
				m.snippet.FocusedStyle.Text = focusedStyle

				m.blur([]string{"description", "promptPass"})

			} else if 1 == m.focusIndex {
				cmds[1] = m.description.Focus()
				m.description.TextStyle = focusedStyle
				m.description.Cursor.Style = focusedStyle

				m.blur([]string{"snippet", "promptPass"})

			} else if 2 == m.focusIndex {
				cmds[2] = m.promptPass.Focus()
				m.promptPass.TextStyle = focusedStyle
				m.promptPass.Cursor.Style = focusedStyle

				m.blur([]string{"snippet", "description"})

			} else if 3 == m.focusIndex {
				m.blur([]string{"snippet", "description", "promptPass"})

			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *snippetModel) blur(inputKeys []string) {

	for _, inputKey := range inputKeys {
		if inputKey == "snippet" {
			m.snippet.Blur()
			m.snippet.Cursor.Style = noStyle
			m.snippet.Cursor.TextStyle = noStyle
		}
		if inputKey == "description" {
			m.description.Blur()
			m.description.TextStyle = noStyle
			m.description.Cursor.Style = noStyle
		}
		if inputKey == "promptPass" {
			m.promptPass.Blur()
			m.promptPass.TextStyle = noStyle
			m.promptPass.Cursor.Style = noStyle
		}
	}

}

func (m *snippetModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.snippet, cmd = m.snippet.Update(msg)
	snippetInfo.Command = m.snippet.Value()
	cmds = append(cmds, cmd)
	m.description, cmd = m.description.Update(msg)
	snippetInfo.Description = m.description.Value()
	cmds = append(cmds, cmd)
	m.promptPass, cmd = m.promptPass.Update(msg)
	snippetInfo.PromptPass = m.promptPass.Value()
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m snippetModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(m.snippet.View())
	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString(m.description.View())
	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString(m.promptPass.View())
	b.WriteRune('\n')

	button := &blurredButton
	if m.focusIndex == 3 {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("[ `tab` to switch, `ctrl+c` to cancel ]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ `enter` on button to save ]"))
	b.WriteRune('\n')

	return b.String()
}

func NewSnippet() *SnippetInfo {

	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	if exitCode == 0 {
		os.Exit(exitCode)
	}

	return &snippetInfo
}

package dialog

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
			Render(" Run ")
	blurredButton = blurredStyle.Copy().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"}).
			Render(" Run ")
	paramIndex = map[string]int{}
	exitCode   = -1
)

type paramFormModel struct {
	focusIndex     int
	commandSnippet textarea.Model
	paramInputs    []textinput.Model
}

func initialModel(params map[string]string, command string) paramFormModel {
	m := paramFormModel{
		paramInputs: make([]textinput.Model, 0),
	}

	ta := textarea.New()
	ta.Placeholder = "Snippet"
	ta.SetValue(command)
	ta.Cursor.Style = noStyle
	ta.Cursor.TextStyle = noStyle
	ta.ShowLineNumbers = false
	ta.CharLimit = 500
	m.commandSnippet = ta

	var t textinput.Model
	idx := 0
	for key, val := range params {

		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.Prompt = ""
		t.PromptStyle = inputStyle
		t.Placeholder = key
		t.CharLimit = 200
		if idx == 0 {
			t.Focus()
		}
		t.SetValue(val)

		m.paramInputs = append(m.paramInputs, t)
		paramIndex[key] = idx
		idx++
	}

	return m
}

func (m paramFormModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m paramFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			exitCode = 0
			return m, tea.Quit

		case "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.paramInputs) {
				return m, tea.Quit
			}

		case "tab", "shift+tab":
			s := msg.String()
			if s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.paramInputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.paramInputs)
			}

			cmds := make([]tea.Cmd, len(m.paramInputs))
			for i := 0; i <= len(m.paramInputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.paramInputs[i].Focus()
					m.paramInputs[i].TextStyle = focusedStyle
					continue
				}
				m.paramInputs[i].Blur()
				m.paramInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *paramFormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.paramInputs))

	for i := range m.paramInputs {
		m.paramInputs[i], cmds[i] = m.paramInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m paramFormModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(m.commandSnippet.View())
	b.WriteRune('\n')

	for i := range m.paramInputs {
		b.WriteRune('\n')
		b.WriteString(blurredStyle.Render(m.paramInputs[i].Placeholder))
		b.WriteRune('\n')
		b.WriteString(m.paramInputs[i].View())
		if i < len(m.paramInputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.paramInputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("[ `tab` to switch, `ctrl+c` to cancel ]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ `enter` on button to save ]"))
	b.WriteRune('\n')

	return b.String()
}

func GenerateParamsForm(params map[string]string, command string) {
	m := initialModel(params, command)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not create cmd form: %s\n", err)
		os.Exit(1)
	}
	if exitCode == 0 {
		os.Exit(exitCode)
	}
	paramsFilled := map[string]string{}

	for k, v := range paramIndex {
		paramsFilled[k] = m.paramInputs[v].Value()
	}
	FinalCommand = insertParams(CurrentCommand, paramsFilled)

}

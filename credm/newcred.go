package credm

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/techierishi/pal/util"
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

	exitCode = -1
)

type passModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel() passModel {
	m := passModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.Prompt = ""
		t.PromptStyle = inputStyle
		t.CharLimit = 200

		switch i {
		case 0:
			t.Placeholder = "Application"
			t.Focus()
			t.PromptStyle = inputStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 500

		case 1:
			t.Placeholder = "Username"
			t.CharLimit = 64
			t.Prompt = ""
			t.PromptStyle = inputStyle
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m passModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m passModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			exitCode = 0
			return m, tea.Quit

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *passModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m passModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteRune('\n')
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("[ `tab` to switch, `ctrl+c` to cancel ]"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("[ `enter` on button to save ]"))
	b.WriteRune('\n')

	return b.String()
}

func NewCred() *CredInfo {
	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	if exitCode == 0 {
		os.Exit(exitCode)
	}
	timestamp := util.UnixMilli()

	credInfo := CredInfo{
		Application: m.inputs[0].Value(),
		Username:    m.inputs[1].Value(),
		Password:    m.inputs[2].Value(),
		Timestamp:   timestamp,
	}

	return &credInfo
}

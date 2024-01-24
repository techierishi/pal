package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	errMessageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#750000"}).
			Render
	selectedSearchItem = SearchRowItem{}
	exitCode           = -1
)

type SearchRowItem struct {
	idx         string
	title, desc string
}

type CustomLabel struct {
	SearchTitle    string
	EnterHelpText  string
	DeleteHelpText string
}

func (cl *CustomLabel) GetEnterHelpText() string {
	if cl.EnterHelpText == "" {
		return "copy to clipboard"
	}
	return cl.EnterHelpText
}

func (cl *CustomLabel) GetDeleteHelpText() string {
	if cl.DeleteHelpText == "" {
		return "delete"
	}
	return cl.DeleteHelpText
}

func NewSearchRowItem(title string, idx string) SearchRowItem {
	return SearchRowItem{
		title: title,
		idx:   idx,
	}
}

func (i SearchRowItem) Title() string       { return i.title }
func (i SearchRowItem) Index() string       { return i.idx }
func (i SearchRowItem) Description() string { return i.desc }
func (i SearchRowItem) FilterValue() string { return i.title }

type searchListModel struct {
	list list.Model
}

func (m searchListModel) Init() tea.Cmd {
	return nil
}

func (m searchListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" || msg.String() == "q" {
			exitCode = 0
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m searchListModel) View() string {
	return docStyle.Render(m.list.View())
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}

	d.SetHeight(2)
	d.SetSpacing(1)

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		selectedItem := SearchRowItem{}

		if i, ok := m.SelectedItem().(SearchRowItem); ok {
			selectedItem.title = i.Title()
			selectedItem.idx = i.Index()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				selectedSearchItem = selectedItem
				return tea.Quit

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted selected item... "))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newDelegateKeyMap(helpText CustomLabel) *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", helpText.GetEnterHelpText()),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", helpText.GetDeleteHelpText()),
		),
	}
}

func SearchUI(customLabel CustomLabel, searchList []list.Item) (SearchRowItem, error) {

	var (
		delegateKeys = newDelegateKeyMap(customLabel)
	)

	delegate := newItemDelegate(delegateKeys)
	searchItemList := list.New(searchList, delegate, 0, 0)
	m := searchListModel{list: searchItemList}
	m.list.Title = customLabel.SearchTitle

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running search ui:", err)
		os.Exit(1)
	}

	if exitCode == 0 {
		os.Exit(exitCode)
	}

	return selectedSearchItem, nil
}

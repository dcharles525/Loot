package ui

import (
	"log"
	"loot/data"
	"loot/structs"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	nav structs.Mode = iota
	edit
	createTitle
	createCommand
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	List      list.Model
	TitleIn   textinput.Model
	CommandIn textinput.Model
	Mode      structs.Mode
	Width     int
	Height    int
}

func initialModel() Model {
	items, _ := data.GetAllItems()
	titleIn := textinput.New()
	titleIn.Prompt = "$ "
	titleIn.Placeholder = "Command Name..."
	titleIn.CharLimit = 100
	titleIn.Width = 50

	commandIn := textinput.New()
	commandIn.Prompt = "$ "
	commandIn.Placeholder = "Command..."
	commandIn.CharLimit = 500
	commandIn.Width = 50

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Commands"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			structs.Keymaps.Create,
			structs.Keymaps.Delete,
			structs.Keymaps.Back,
		}
	}

	return Model{
		List:      l,
		TitleIn:   titleIn,
		CommandIn: commandIn,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println(msg)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.TitleIn.Focused() {
			if m.Mode == createTitle {
				if key.Matches(msg, structs.Keymaps.Enter) {
					m.Mode = createCommand
					m.TitleIn.Blur()
					m.CommandIn.Focus()
					cmds = append(cmds, textinput.Blink)
				}
			}
			var cmd tea.Cmd
			m.TitleIn, cmd = m.TitleIn.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.CommandIn.Focused() {
			if key.Matches(msg, structs.Keymaps.Enter) {
				cmds = append(cmds, data.CreateCommand(
					m.TitleIn.Value(),
					m.CommandIn.Value(),
				))

				m.TitleIn.SetValue("")
				m.CommandIn.SetValue("")
				m.CommandIn.Blur()
				m.Mode = nav
			}

			var cmd tea.Cmd
			m.CommandIn, cmd = m.CommandIn.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, structs.Keymaps.Create):
				m.Mode = createTitle
				m.TitleIn.Focus()
				cmds = append(cmds, textinput.Blink)
			case key.Matches(msg, structs.Keymaps.Delete):
				items := m.List.Items()
				if len(items) > 0 {
					cmds = append(cmds, data.DeleteCommand(m.getActiveProjectID()))
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	case structs.UpdateCommandMsg:
		m.Mode = nav
	}

	// Refresh list if not filtering
	if m.List.FilterState() != list.Filtering {
		items, _ := data.GetAllItems()
		m.List.SetItems(items)
	}

	if msg != "q" && (m.Mode != createTitle && m.Mode != createCommand) {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.TitleIn.Focused() {
		return docStyle.Render(m.List.View() + "\n" + m.TitleIn.View())
	}
	if m.CommandIn.Focused() {
		return docStyle.Render(m.List.View() + "\n" + m.CommandIn.View())
	}
	return docStyle.Render(m.List.View())
}

func (m Model) getActiveProjectID() string {
	items := m.List.Items()
	if len(items) == 0 {
		return ""
	}
	activeItem := items[m.List.Index()]
	return activeItem.(structs.Item).CommandTitle
}

func Run() error {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

package main

import (
	"fmt"
	"log"
	"os"

	structs "loot/structs"

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

type model struct {
	list list.Model
	commandTitleInput	textinput.Model
	commandInput textinput.Model
	mode structs.Mode
	width int
	height int
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println(msg)
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.commandTitleInput.Focused(){
			if m.mode == createTitle {
				if key.Matches(msg, structs.Keymaps.Enter) {
					m.mode = createCommand
					m.commandTitleInput.Blur()
					m.commandInput.Focus()
					
					var cmd tea.Cmd
					cmd = textinput.Blink
					cmds = append(cmds, cmd)
				}
				
			}
			var cmd tea.Cmd
			m.commandTitleInput, cmd = m.commandTitleInput.Update(msg)
			cmds = append(cmds, cmd)
		}else if m.commandInput.Focused() {
			if key.Matches(msg, structs.Keymaps.Enter) {	
				cmds = append(cmds, CreateCommand(
					m.commandTitleInput.Value(), 
					m.commandInput.Value(),
				))
				
				m.commandTitleInput.SetValue("")
				m.commandInput.SetValue("")
				m.commandInput.Blur()
				m.mode = nav
			}

			var cmd tea.Cmd
			m.commandInput, cmd = m.commandInput.Update(msg)
			cmds = append(cmds, cmd)
		}else{
			switch {
			case key.Matches(msg, structs.Keymaps.Create):
				m.mode = createTitle
				m.commandTitleInput.Focus()
				
				var cmd tea.Cmd
				cmd = textinput.Blink
				cmds = append(cmds, cmd)
			case key.Matches(msg, structs.Keymaps.Delete):
				items := m.list.Items()
				
				var cmd tea.Cmd
				if len(items) > 0 {
					cmd = DeleteCommand(m.getActiveProjectID())
				}
				cmds = append(cmds, cmd)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case structs.UpdateCommandMsg:
		m.mode = nav
	}
	
	if m.list.FilterState() != list.Filtering {
		items := []list.Item{}
		items = append(items, GetAllItems()...)
		m.list.SetItems(items)
	}


	if msg != "q" && (m.mode != createTitle &&  m.mode != createCommand) {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.commandTitleInput.Focused() {
		return docStyle.Render(m.list.View() + "\n" + m.commandTitleInput.View())
	}
	if m.commandInput.Focused() {
		return docStyle.Render(m.list.View() + "\n" + m.commandInput.View())
	}
	return docStyle.Render(m.list.View())
}

func (m model) getActiveProjectID() string {
	items := m.list.Items()
	activeItem := items[m.list.Index()]
	return activeItem.(structs.Item).CommandTitle
}

func main() {
	file, err := os.OpenFile("loot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	
	items := []list.Item{}
	items = append(items, GetAllItems()...)
	
	commandTitleInput := textinput.New()
	commandTitleInput.Prompt = "$ "
	commandTitleInput.Placeholder = "Command Name..."
	commandTitleInput.CharLimit = 100
	commandTitleInput.Width = 50

	commandInput := textinput.New()
	commandInput.Prompt = "$ "
	commandInput.Placeholder = "Command..."
	commandInput.CharLimit = 500
	commandInput.Width = 50

	m := model{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
		commandInput: commandInput,
		commandTitleInput: commandTitleInput,
	}
	m.list.Title = "Commands"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			structs.Keymaps.Create,
			structs.Keymaps.Delete,
			structs.Keymaps.Back,
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

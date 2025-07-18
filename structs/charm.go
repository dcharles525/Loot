package structs

import "github.com/charmbracelet/bubbles/key"

type Mode int

const (
	Nav Mode = iota
	Edit
	CreateTitle
	CreateCommand
)

type UpdateCommandMsg struct{}

type Keymap struct {
	Create key.Binding
	Enter  key.Binding
	Rename key.Binding
	Delete key.Binding
	Back   key.Binding
	Quit   key.Binding
}

var Keymaps = Keymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}

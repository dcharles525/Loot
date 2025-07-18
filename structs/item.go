package structs

type Item struct {
	CommandTitle, CommandDesc string
}

func (i Item) Title() string       { return i.CommandTitle }
func (i Item) Description() string { return i.CommandDesc }
func (i Item) FilterValue() string { return i.CommandTitle }